package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"go-purple/purple"
)

type PurpleServer struct {
	pc *purple.PurpleCore
	mh *httprouter.Router
}

func newPurpleServer() *PurpleServer {
	this := &PurpleServer{}
	return this
}

func (this *PurpleServer) serve() {
	go func() { this.setupPurpleCore() }()

	this.setupRouter()

	time.Sleep(100 * time.Millisecond)
	port := 6363
	log.Printf("started: *:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), this.mh)
	log.Println("stopped", err)
}

func (this *PurpleServer) setupPurpleCore() {
	this.pc = purple.NewPurpleCore()

	this.pc.MainLoop()
	log.Println("started purple core")
}

func (this *PurpleServer) setupRouter() {
	dummy := httprouter.New()
	dummy.GET("/", this.indexPage)

	this.mh = dummy
}

func (this *PurpleServer) indexPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acs := this.pc.AccountsGetAll()
	for _, ac := range acs {
		conn := ac.GetConnection()
		plugInfo := conn.GetPrplInfo()
		log.Println(ac.GetUserName(), ac.GetAlias(), plugInfo.Name, plugInfo.Id)

		buf := bytes.NewBuffer(nil)
		buf.WriteString(ac.GetUserName())
		buf.WriteString(ac.GetAlias())
		buf.WriteString(plugInfo.Name)
		buf.WriteString(plugInfo.Id)
		buf.WriteTo(w)
	}
}
