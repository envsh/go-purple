/*
 解析用户信息数据
 用户信息数据结构
*/
package wechat

import (
	"fmt"
	"log"

	"github.com/bitly/go-simplejson"
)

type User struct {
	UserName  string
	NickName  string
	Signature string
	// ContactFlag int // 好像是联系活跃度

	Members []*User
}

func parseContact(contacto *simplejson.Json) *User {
	u := &User{}
	u.UserName = contacto.Get("UserName").MustString()
	u.NickName = contacto.Get("NickName").MustString()
	// log.Println(u)

	memcnt := contacto.Get("MemberCount").MustInt()
	if memcnt > 0 {
		u.Members = make([]*User, memcnt)
		for idx := 0; idx < memcnt; idx++ {
			memo := contacto.Get("MemberList").GetIndex(idx)
			mu := parseContact(memo)
			u.Members[idx] = mu
		}
	}

	return u
}

func parseContact2(contacto *simplejson.Json) *User {
	u := &User{}
	u.UserName = contacto.Get("UserName").MustString()
	u.NickName = contacto.Get("NickName").MustString()
	// log.Println(u)

	u.Members = make([]*User, 0)
	NewParser2(contacto).Each("Member", func(itemo *simplejson.Json) {
		mu := parseContact(itemo)
		u.Members = append(u.Members, mu)
	})

	return u
}

func parseWXInitData(data string) (users []*User) {
	jso, err := simplejson.NewJson([]byte(data))
	if err != nil {
		log.Println(err)
		return
	}
	retv := jso.GetPath("BaseResponse", "Ret").MustInt()
	if retv != 0 {
		log.Println("Invalid resonse:", retv)
		return
	}

	cntv := jso.GetPath("Count").MustInt()
	log.Println("parsering contact:", cntv)
	users = make([]*User, cntv)
	for idx := 0; idx < cntv; idx++ {
		contacto := jso.Get("ContactList").GetIndex(idx)
		u := parseContact(contacto)
		users[idx] = u
	}

	// MPArticleList
	return
}

func parseContactData(data string) (users []*User) {
	jso, err := simplejson.NewJson([]byte(data))
	if err != nil {
		log.Println(err)
		return
	}
	retv := jso.GetPath("BaseResponse", "Ret").MustInt()
	if retv != 0 {
		log.Println("Invalid resonse:", retv)
		return
	}

	cntv := jso.GetPath("MemberCount").MustInt()
	log.Println("parsering member:", cntv)
	users = make([]*User, cntv)
	for idx := 0; idx < cntv; idx++ {
		contacto := jso.Get("MemberList").GetIndex(idx)
		u := parseContact(contacto)
		users[idx] = u
	}

	return
}

func parseContactData2(data string) (users []*User) {
	p := NewParser(data)
	if !p.RetOK() {
		return
	}

	users = make([]*User, 0)
	p.Each("Member", func(itemo *simplejson.Json) {
		u := parseContact(itemo)
		users = append(users, u)
	})

	return
}

// 通用版本
type Parser struct {
	data string
	jso  *simplejson.Json
}

func NewParser(data string) *Parser {
	this := &Parser{}
	this.data = data

	this.init()
	return this
}

func NewParser2(jso *simplejson.Json) *Parser {
	this := &Parser{}
	this.jso = jso

	return this
}

func (this *Parser) init() {
	var err error
	this.jso, err = simplejson.NewJson([]byte(this.data))
	if err != nil {
		log.Println(err)
	}
}

func (this *Parser) RetOK() bool {
	if this.jso == nil {
		return false
	}
	if len(this.data) == 0 {
		return true
	}
	retv := this.jso.GetPath("BaseResponse", "Ret").MustInt()
	if retv != 0 {
		log.Println("Invalid resonse:", retv)
		return false
	}
	return true
}

// items 包括 AddMsg, Member, Contact, MPArticle
func (this *Parser) Each(item string, functor func(itemo *simplejson.Json)) {
	cntkey := fmt.Sprintf("%sCount", item)
	lstkey := fmt.Sprintf("%sList", item)

	if this.jso.Get(cntkey) != nil {
		cntv := this.jso.Get(cntkey).MustInt()
		for idx := 0; idx < cntv; idx++ {
			functor(this.jso.Get(lstkey).GetIndex(idx))
		}
	}
}
