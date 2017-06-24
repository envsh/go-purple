package main

import (
	"log"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/chrislusf/glow/flow"
	"github.com/comail/colog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nats-io/go-nats"
	"github.com/seiflotfy/cuckoofilter"
	"github.com/willf/bloom"
)

type AnaContext struct {
	fctx    *flow.FlowContext
	mbcli   *MsgBusClient
	dbh     *gorm.DB
	ana     *Analizer
	blmflt  *bloom.BloomFilter
	cuckflt *cuckoofilter.CuckooFilter
	fch     chan *EventR
}

func NewAnaContext() *AnaContext {
	this := &AnaContext{}

	this.fctx = flow.New()
	this.blmflt = bloom.New(100000, 5)
	this.cuckflt = cuckoofilter.NewDefaultCuckooFilter()
	this.fch = make(chan *EventR, 0)
	this.ana = NewAnalizer()

	return this
}

var actx = NewAnaContext()

type Analizer struct {
	fctx *flow.FlowContext
}

func NewAnalizer() *Analizer {
	this := &Analizer{}
	this.fctx = flow.New()
	return this
}
func (this *Analizer) run() {
	fds := this.fctx.Channel(actx.fch)
	fds.Filter(func(line *EventR) bool {
		log.Println(line)
		return true
	}).Map(func(line *EventR, ch chan rune) {
		log.Println()
		for _, r := range line.Message {
			ch <- r
		}
	}).Map(func(r rune) (int, int8) {
		return 1, 2
	}).Reduce(func(x flow.KeyValue, y flow.KeyValue) flow.KeyValue {
		log.Println(x, y)
		return flow.KeyValue{Key: 5, Value: 6}
	})

	fds.Run()
}

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
	if false { // use too much memory, about 1G
		nlp := NewSnowNLP(string(m.Data))
		log.Println(nlp.Sentiments(), nlp.Words())
	}

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

	actx.fch <- evtrec

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
	log.Printf("%+v\n", db.DB().Driver())
	dbh = db
	actx.dbh = db
	defer db.Close()

	db.AutoMigrate(&EventR{})
	go actx.ana.run()

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
