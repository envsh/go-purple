package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (this *PurpleServer) infoPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acs := this.pc.AccountsGetAll()
	for _, ac := range acs {
		conn := ac.GetConnection()
		plugInfo := conn.GetPrplInfo()
		log.Println(ac.GetUserName(), ac.GetAlias(), plugInfo.Name, plugInfo.Id)
	}
}
