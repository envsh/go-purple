package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"yobot/purple"
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
	this.setupPurpleCore()

	this.setupRouter()
	log.Println("started")
	err := http.ListenAndServe(":6363", this.mh)
	log.Println("stopped", err)
}

func (this *PurpleServer) setupPurpleCore() {
	this.pc = purple.NewPurpleCore()
	this.pc.InitUi()
	this.pc.InitPurple()

	go this.pc.Loop()
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
