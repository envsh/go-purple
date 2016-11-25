/*
 解析用户信息数据
 用户信息数据结构
*/
package wechat

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

type User struct {
	UserName  string
	NickName  string
	Signature string
	// ContactFlag int // 好像是联系活跃度

	Members []*User
}

func (this *User) Equal(u *User) bool {
	return u != nil && u.UserName == this.UserName
}

func (this *User) IsGroup() bool {
	return strings.HasPrefix(this.UserName, "@@")
}

func ParseContact(contacto *simplejson.Json) *User {
	u := &User{}
	u.UserName = contacto.Get("UserName").MustString()
	u.NickName = contacto.Get("NickName").MustString()
	// log.Println(u)

	u.Members = make([]*User, 0)
	NewParser2(contacto).Each("Member", func(itemo *simplejson.Json) {
		mu := ParseContact(itemo)
		u.Members = append(u.Members, mu)
	})

	return u
}

// usernames seperated by ","
func parseChatSet(set *simplejson.Json) (cset []*User) {
	cset = make([]*User, 0)
	if seto, ok := set.CheckGet("ChatSet"); ok {
		for _, name := range strings.Split(seto.MustString(), ",") {
			u := &User{UserName: name}
			cset = append(cset, u)
		}
	}
	return
}

type MPArticle struct{}

func ParseWXInitData(data string) (users []*User) {
	p := NewParser(data)
	if !p.RetOK() {
		return
	}

	users = make([]*User, 0)
	p.Each("Contact", func(itemo *simplejson.Json) {
		u := ParseContact(itemo)
		users = append(users, u)
	})
	log.Println("parsered contacts:", len(users))

	// MPArticleList
	mpas := make([]*MPArticle, 0)
	p.Each("MPArticle", func(itemo *simplejson.Json) {
		mpa := &MPArticle{}
		mpas = append(mpas, mpa)
	})
	log.Println("parsered mpas:", len(mpas))

	// MPSubscribeMsgList
	mpsubs := make([]*MPArticle, 0)
	p.Each("MPSubscribeMsg", func(itemo *simplejson.Json) {
		mpsub := &MPArticle{}
		mpsubs = append(mpsubs, mpsub)
	})
	log.Println("parsered mpsubs:", len(mpas))

	// chatset
	cset := parseChatSet(p.jso)
	log.Println("ChatSet", len(cset), cset)

	return
}

func ParseContactData(data string) (users []*User) {
	p := NewParser(data)
	if !p.RetOK() {
		return
	}

	users = make([]*User, 0)
	p.Each("Member", func(itemo *simplejson.Json) {
		u := ParseContact(itemo)
		users = append(users, u)
	})
	log.Println("parsered members:", len(users))

	return
}

// 通用版本
// TODO 支持非JSON格式响应数据的解析
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

	var cntv int
	if jv, ok := this.jso.CheckGet(cntkey); ok {
		cntv = jv.MustInt()
	} else if jv, ok := this.jso.CheckGet("Count"); ok {
		cntv = jv.MustInt()
	} else {
		// log.Printf("Count and %sCount both not found.", item)
	}

	for idx := 0; idx < cntv; idx++ {
		functor(this.jso.Get(lstkey).GetIndex(idx))
	}
}

// golang的静态函数好像是可以的，只要不使用this
// 用法：Parser{}.parseUUID()
func (this Parser) parseUUID(str string) (code int, uuid string) {
	exp := regexp.MustCompile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "([\w\-=]+)";`)
	mats := exp.FindAllStringSubmatch(str, -1)
	if len(mats) > 0 {
		code, _ = strconv.Atoi(mats[0][1])
		uuid = mats[0][2]
	}
	return
}

func (this Parser) parseScan(str string) (code int) {
	exp := regexp.MustCompile(`window.code=(\d+);`)
	mats := exp.FindAllStringSubmatch(str, -1)
	if len(mats) > 0 {
		code, _ = strconv.Atoi(mats[0][1])
	}
	return
}

func (this Parser) parseTicket(str string) (ret int, skey string, wxsid string,
	wxuin string, pass_ticket string) {
	// `<error><ret>0</ret><message>OK</message><skey>@crypt_3ea2fe08_723d1e1bd7b4171657b58c6d2849b367</skey><wxsid>9qxNHGgi9VP4/Tx6</wxsid><wxuin>979270107</wxuin><pass_ticket>%2BEdqKi12tfvM8ZZTdNeh4GLO9LFfwKLQRpqWk8LRYVWFkDE6%2FZJJXurz79ARX%2FIT</pass_ticket><isgrayscale>1</isgrayscale></error>`
	exp := regexp.MustCompile(`<error><ret>(\d+)</ret><message>.*</message><skey>(.+)</skey><wxsid>(.+)</wxsid><wxuin>(\d+)</wxuin><pass_ticket>(.+)</pass_ticket><isgrayscale>1</isgrayscale></error>`)
	mats := exp.FindAllStringSubmatch(str, -1)
	if len(mats) > 0 {
		ret, _ = strconv.Atoi(mats[0][1])
		skey = mats[0][2]
		wxsid = mats[0][3]
		wxuin = mats[0][4]
		pass_ticket = mats[0][5]
	}

	return
}

func (this Parser) parseSyncCheck(str string) (code int, selector int) {
	exp := regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
	mats := exp.FindAllStringSubmatch(str, -1)
	if len(mats) > 0 {
		code, _ = strconv.Atoi(mats[0][1])
		selector, _ = strconv.Atoi(mats[0][2])
	} else {
		code = -1
		selector = -1
	}
	return
}
