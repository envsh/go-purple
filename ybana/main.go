package main

import (
	"log"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/comail/colog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nats-io/go-nats"
	"github.com/seiflotfy/cuckoofilter"
	"github.com/willf/bloom"
)

// publish all messages to nats message bus
type MsgBusClient struct {
	nc *nats.Conn
}

func newMsgBusClient() *MsgBusClient {
	this := &MsgBusClient{}
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Println(err)
	}

	this.nc = nc
	return this
}

var blmflt = bloom.New(100000, 5)
var cuckflt = cuckoofilter.NewDefaultCuckooFilter()

var dbh *gorm.DB

type EventR struct {
	// gorm.Model
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Proto   string
	Nick    string
	Ident   string
	EType   string
	Message string
}

func messageHandler(m *nats.Msg) {
	// fmt.Printf("Received a message: %s\n", string(m.Data))
	jso, err := simplejson.NewJson(m.Data)
	if err != nil {
		log.Println(err)
		return
	}
	// TODO convert back to Event?

	switch jso.Get("Proto").MustString() {
	case "table":
		return
	}
	bfok := blmflt.TestAndAdd(m.Data)
	// log.Println(blmflt.TestAndAdd(m.Data))
	cfok := cuckflt.InsertUnique(m.Data)
	// log.Println(cuckflt.InsertUnique(m.Data))
	if bfok != !cfok {
		log.Println(bfok, cfok) // filter结果不一致
	}
	if bfok == true || cfok == false {
		return // filtered
	}

	log.Printf("Received a message: %s\n", string(m.Data))
	nlp := NewSnowNLP(string(m.Data))
	log.Println(nlp.Sentiments(), nlp.Words())

	evtrec := &EventR{}
	evtrec.EType = jso.Get("EType").MustString()
	evtrec.Proto = jso.Get("Proto").MustString()
	evtrec.Ident = jso.Get("Ident").MustString()

	switch jso.Get("Proto").MustString() {
	case "tox":
		evtrec.Nick = ""
		evtrec.Message = jso.Get("Args").GetIndex(0).MustString()
	case "irc":
		evtrec.Nick = jso.Get("Args").GetIndex(0).MustString()
		evtrec.Message = jso.Get("Args").GetIndex(2).MustString()
	}

	dbh.Create(evtrec)
	if dbh.Error != nil {
		log.Println(dbh.Error)
	}
}

func main() {
	db, err := gorm.Open("sqlite3", "ybana.db")
	if err != nil {
		log.Println(err)
		return
	}
	dbh = db
	defer db.Close()

	db.AutoMigrate(&EventR{})

	mbc := newMsgBusClient()
	sc, err := mbc.nc.Subscribe("yobotmsg", messageHandler)
	if err != nil {
		log.Println(err, sc.Type())
	}
	select {}
}

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	colog.Register()
}
