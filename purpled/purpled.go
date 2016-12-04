package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-purple/purple"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/flosch/pongo2.v3"
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
	log.Println(purple.UserDir())
	var userDir = os.Getenv("HOME") + "/.gopurple"
	purple.UtilSetUserDir(userDir)
	log.Println(purple.UserDir())
	this.pc = purple.NewPurpleCore()

	this.pc.MainLoop()
	log.Println("started purple core")
}

func (this *PurpleServer) setupRouter() {
	dummy := httprouter.New()
	dummy.GET("/dist/:distname", this.distPage)
	dummy.GET("/", this.indexPage)

	this.mh = dummy
}

func AssetHmtl(name string) string { return string(MustAsset("html/" + name)) }
func AssetCss(name string) string  { return string(MustAsset("css/" + name)) }
func AssetJs(name string) string   { return string(MustAsset("js/" + name)) }

func (this *PurpleServer) indexPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pctx := pongo2.Context{}
	accounts := make([]map[string]string, 0)

	acs := this.pc.AccountsGetAll()
	log.Println(len(acs))
	for _, ac := range acs {
		conn := ac.GetConnection()
		account := make(map[string]string, 0)
		account["UserName"] = ac.GetUserName()
		account["Alias"] = ac.GetAlias()
		account["ProtocolId"] = ac.GetProtocolId()
		account["ProtocolName"] = ac.GetProtocolName()
		accounts = append(accounts, account)

		if conn == nil {
			log.Println(ac.GetUserName(), ac.GetAlias(), ac.GetProtocolId(), ac.GetProtocolName())
			continue
		}
		plugInfo := conn.GetPrplInfo()
		log.Println(ac.GetUserName(), ac.GetAlias(), plugInfo.Name, plugInfo.Id)

		/*
			buf := bytes.NewBuffer(nil)
			buf.WriteString(ac.GetUserName())
			buf.WriteString(ac.GetAlias())
			buf.WriteString(plugInfo.Name)
			buf.WriteString(plugInfo.Id)
			buf.WriteTo(w)
		*/
	}

	pctx["accounts"] = accounts

	tplcc := AssetHmtl("index.html")
	tpl, err := pongo2.FromString(tplcc)
	if err != nil {
		log.Println(err)
		buf := bytes.NewBuffer(nil)
		buf.WriteString(err.Error())
		buf.WriteTo(w)
	} else {
		tpl.ExecuteWriter(pctx, w)
	}
}

func (this *PurpleServer) distPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(r.URL.Path, ps)
	p0 := ps[0]
	w.Write([]byte(AssetHmtl("dist/" + p0.Value)))
}
