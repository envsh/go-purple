package wechat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/bitly/go-simplejson"
	"github.com/levigross/grequests"
)

type pollState struct {
	qruuid  string
	qrpic   []byte
	logined bool

	redirUrl    string
	urlBase     string
	pushUrlBase string

	wxdevid      string
	wxuin        string // in cookie
	wxsid        string // in cookie
	wxDataTicket string // in cookie
	wxSKeyOld    string
	wxPassTicket string
	wxuvid       string // in cookie
	wxAuthTicket string // in cookie

	wxSyncKey        *simplejson.Json
	wxSKey           string
	wxInitRawData    string
	wxContactRawData string
}

// 用于保持连接和接收消息，
// 如果要发送消息，需要在登陆成功后在另外的线程使用该rses发送请求。
// 需要考虑的是发送请求是否需要队列，还是并发的呢？
type longPoll struct {
	eqch     chan<- *Event
	rses     *grequests.Session
	rops     *grequests.RequestOptions
	reqState int

	//
	state pollState

	// persistent
	cookies []*http.Cookie
}

func newLongPoll(eqch chan<- *Event) *longPoll {
	this := &longPoll{}
	this.eqch = eqch
	this.state.wxdevid = "e669767113868187"
	this.rses = grequests.NewSession(nil)
	this.rops = &grequests.RequestOptions{}
	this.rops.Headers = map[string]string{
		"Referer": "https://wx2.qq.com/?lang=en_US",
	}
	this.rops.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.100 Safari/537.36 Vivaldi/1.5.658.31"

	this.rops.RequestTimeout = 5 * time.Second

	this.loadCookies()
	if this.cookies != nil {
		this.rops.Cookies = this.cookies
	}

	return this
}

func (this *longPoll) start() {
	go this.run()
}

const (
	REQ_NONE int = iota
	REQ_LOGIN
	REQ_QRCODE
	REQ_WAIT_SCAN
	REQ_REDIR_LOGIN
	REQ_WXINIT
	REQ_CONTACT
	REQ_SYNC_CHECK
	REQ_WEB_SYNC
	REQ_END
)

// blocked
func (this *longPoll) run() {
	if this.cookies != nil {
		log.Println("check cookies validation...")
		if !this.checkLoadedCookiesByPingOrNoop() {
			log.Println("cookies invalid clean up...")
			this.resetState()
		}
	}

	log.Println("polling...")
	stopped := false
	if this.cookies != nil {
		this.reqState = REQ_SYNC_CHECK
		this.reqState = REQ_CONTACT
	} else {
		this.reqState = REQ_LOGIN
	}

	for !stopped {
		switch this.reqState {
		case REQ_LOGIN:
			this.jslogin()
		case REQ_QRCODE:
			this.getqrcode()
		case REQ_WAIT_SCAN:
			this.pollScan()
		case REQ_REDIR_LOGIN:
			this.redirLogin()
		case REQ_WXINIT:
			this.wxInit()
		case REQ_CONTACT:
			this.getContact()
		case REQ_SYNC_CHECK:
			this.syncCheck()
		case REQ_WEB_SYNC:
			this.webSync()
		case REQ_END:
			this.eqch <- newEvent(EVT_LOGOUT, []string{})
			stopped = true
		}
	}

	log.Println("run end")
}

func parseuuid(str string) (code int, uuid string) {
	exp := regexp.MustCompile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "([\w\-=]+)";`)
	mats := exp.FindAllStringSubmatch(str, -1)
	if len(mats) > 0 {
		code, _ = strconv.Atoi(mats[0][1])
		uuid = mats[0][2]
	}
	return
}

func (this *longPoll) saveContent(name string, bcc []byte, resp *grequests.Response, url string) {
	err := ioutil.WriteFile(name, bcc, 0644)
	if err != nil {
		log.Println(err)
	}
}

var cookies_file = "./cookies.json"

// TODO 检测加载的状态数据是否还能够使用
func (this *longPoll) loadCookies() {
	sck, err := ioutil.ReadFile("cookies.txt")
	if err != nil {
		log.Println(err)
	} else {
		jck, err := simplejson.NewJson(sck)
		if err != nil {
			log.Println(err)
		} else {
			ckarr := jck.Get("cookies").MustArray()
			if len(ckarr) == 0 {
				log.Println("Invalid json node")
				return
			}
			this.state.qruuid = jck.Get("qruuid").MustString()
			this.state.wxSKey = jck.Get("wxskey").MustString()
			this.state.wxPassTicket = jck.Get("pass_ticket").MustString()
			this.state.redirUrl = jck.Get("redir_url").MustString()
			this.state.urlBase = jck.Get("urlBase").MustString()
			this.state.pushUrlBase = jck.Get("pushUrlBase").MustString()
			this.state.wxSyncKey = jck.Get("SyncKey")

			this.cookies = make([]*http.Cookie, 0)
			for idx, _ := range ckarr {
				hck := &http.Cookie{}
				bck, err := jck.Get("cookies").GetIndex(idx).MarshalJSON()
				if err != nil {
				}
				err = json.Unmarshal(bck, hck)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(idx, hck)
					this.cookies = append(this.cookies, hck)
					if hck.Name == "wxuin" {
						this.state.wxuin = hck.Value
					} else if hck.Name == "wxsid" {
						this.state.wxsid = hck.Value
					} else if hck.Name == "webwx_data_ticket" {
						this.state.wxDataTicket = hck.Value
					} else if hck.Name == "webwxuvid" {
						this.state.wxuvid = hck.Value
					} else if hck.Name == "webwx_auth_ticket" {
						this.state.wxAuthTicket = hck.Value
					}
				}
			}
		}
	}
}

func (this *longPoll) saveCookies(resp *grequests.Response) {
	/*
		for _, ck := range resp.RawResponse.Cookies() {
			log.Println(ck)
		}
		for hkey, hval := range resp.Header {
			log.Println(hkey, "=", hval)
		}
	*/

	var jck *simplejson.Json
	bck, err := ioutil.ReadFile("cookies.txt")
	if err != nil {
		log.Println(err)
		jck = simplejson.New()
	} else {
		jck, err = simplejson.NewJson(bck)
	}

	if len(resp.RawResponse.Cookies()) > 0 {
		// 合并cookies：先解析旧数据，再添加到新数据上
		ckarr := jck.Get("cookies").MustArray()
		cookies := make([]*http.Cookie, 0)
		for idx, _ := range ckarr {
			hck := &http.Cookie{}
			bck, err := jck.Get("cookies").GetIndex(idx).MarshalJSON()
			if err != nil {
			}
			err = json.Unmarshal(bck, hck)
			cookies = append(cookies, hck)
		}
		newcookies := resp.RawResponse.Cookies()
		for _, ckold := range cookies {
			found := false
			for _, ck := range newcookies {
				if ckold.Name == ck.Name {
					found = true
					break
				}
			}
			if !found { // refill it
				newcookies = append(newcookies, ckold)
			}
		}
		jck.Set("cookies", newcookies)
	}
	jck.Set("qruuid", this.state.qruuid)
	jck.Set("wxskey", this.state.wxSKey)
	jck.Set("pass_ticket", this.state.wxPassTicket)
	jck.Set("redir_url", this.state.redirUrl)
	jck.Set("urlBase", this.state.urlBase)
	jck.Set("pushUrlBase", this.state.pushUrlBase)
	jck.Set("SyncKey", this.state.wxSyncKey)

	sck, err := jck.Encode()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(sck))
		sck, err = jck.EncodePretty()
		this.saveContent("cookies.txt", sck, nil, "")
	}
}

// 清除状态数据，从扫码开始登陆
func (this *longPoll) resetState() {
	this.cookies = nil
	this.state = pollState{}
}

var ips = strings.Split("101.226.76.164 101.227.160.102 140.206.160.161 140.207.135.104 117.135.169.34 117.144.242.33 203.205.151.221", " ")

// 在加载完数据，做其他操作之前尝试一下是否cookie过期。
// 还需要找一个有效的请求，要求请求速度快，并且需要cookie的。
func (this *longPoll) checkLoadedCookiesByPingOrNoop() bool {
	this.webSync()
	if this.reqState == REQ_END {
		return false
	}
	return true
}

func (this *longPoll) jslogin() {
	url := "https://login.weixin.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx2.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=en_US"
	log.Println(url)
	resp, err := this.rses.Get(url, this.rops)
	if err != nil {
		log.Println(err, url)
	}
	log.Println(resp.StatusCode, resp.Header, resp.Ok)
	bcc := resp.Bytes()
	this.saveContent("jslogin.json", bcc, resp, url)
	defer resp.Close()

	// # parse hcc: window.QRLogin.code = 200; window.QRLogin.uuid = "gYmgd1grLg==";
	code, uuid := parseuuid(resp.String())
	if code != 200 {
		log.Println(resp.String())
		this.reqState = REQ_END
	} else {
		this.state.qruuid = uuid
		this.reqState = REQ_QRCODE
		// this.saveCookies(resp)
		this.eqch <- newEvent(EVT_GOT_UUID, []string{uuid})
	}

}

func (this *longPoll) getqrcode() {
	nsurl := "https://login.weixin.qq.com/qrcode/4ZYgra8RHw=="
	nsurl = "https://login.weixin.qq.com/qrcode/" + this.state.qruuid
	log.Println(nsurl)

	resp, err := this.rses.Get(nsurl, this.rops)
	if err != nil {
		log.Println(err, nsurl)
	}
	log.Println(resp.StatusCode, resp.Header, resp.Ok)
	bcc := resp.Bytes()
	this.saveContent("qrcode.jpg", bcc, resp, nsurl)
	defer resp.Close()
	this.state.qrpic = bcc

	if !resp.Ok {
		this.reqState = REQ_END
	} else {
		this.reqState = REQ_WAIT_SCAN
		// this.saveCookies(resp)
		this.eqch <- newEvent(EVT_GOT_QRCODE, []string{resp.String()})

	}

}

func (this *longPoll) nowTime() int64 {
	return time.Now().Unix()
}

func parsescan(str string) (code int) {
	exp := regexp.MustCompile(`window.code=(\d+);`)
	mats := exp.FindAllStringSubmatch(str, -1)
	if len(mats) > 0 {
		code, _ = strconv.Atoi(mats[0][1])
	}
	return
}

func (this *longPoll) pollScan() {
	nsurl := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=4eDUw9zdPg==&tip=0&r=-1166218796"
	// # v2 url: https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=gfNC8TeiPg==&tip=1&r=-1222670084&lang=en_US
	nsurl = fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=%s&tip=0&r=%d&lang=en_US&_=%d",
		this.state.qruuid, this.nowTime(), this.nowTime())
	log.Println(nsurl)

	resp, err := this.rses.Get(nsurl, this.rops)
	if err != nil {
		log.Println(err, nsurl)
	}
	log.Println(resp.StatusCode, resp.Ok, resp.Header, resp.String())
	bcc := resp.Bytes()
	this.saveContent("pollscan.json", bcc, resp, nsurl)
	defer resp.Close()

	if !resp.Ok {
		this.reqState = REQ_END
	} else {
		/*
					# window.code=408;  # 像是超时
					# window.code=400;  # ??? 难道是会话过期???需要重新获取QR图（已确认，在浏览器中，收到400后刷新了https://wx2.qq.com/
			            # window.code=201;  # 已扫描，未确认
			            # window.code=200;  # 已扫描，已确认登陆
			            # parse hcc, format: window.code=201;
		*/
		code := parsescan(resp.String())

		switch code {
		case 408:
			this.reqState = REQ_WAIT_SCAN // no change
		case 400:
			log.Println("maybe need rerun refresh()...")
			this.reqState = REQ_END
		case 201:
			time.Sleep(2 * time.Second)
		case 200:
			this.state.logined = true
			this.state.redirUrl = strings.Split(resp.String(), "\"")[1]
			this.reqState = REQ_REDIR_LOGIN
		default:
			log.Println("not impled", code)
			this.reqState = REQ_END
		}

		// this.saveCookies(resp)
		this.eqch <- newEvent(EVT_SCAN_DATA, []string{resp.String()})

	}

}

func parseTicket(str string) (ret int, skey string, wxsid string,
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

func (this *longPoll) redirLogin() {
	nsurl := this.state.redirUrl + "&fun=new&version=v2"
	if strings.Contains(nsurl, "wx.qq.com") {
		this.state.urlBase = "https://wx.qq.com"
		this.state.pushUrlBase = "https://webpush.weixin.qq.com"
	} else {
		this.state.urlBase = "https://wx2.qq.com"
		this.state.pushUrlBase = "https://webpush2.weixin.qq.com"
	}
	log.Println(nsurl)

	resp, err := this.rses.Get(nsurl, this.rops)
	if err != nil {
		log.Println(err, nsurl)
	}
	log.Println(resp.StatusCode, resp.Header, resp.Ok)
	bcc := resp.Bytes()
	this.saveContent("redir.html", bcc, resp, nsurl)
	defer resp.Close()

	if !resp.Ok {
		this.reqState = REQ_END
	} else {
		/*
			# parse content: SKey,pass_ticket
			# <error><ret>0</ret><message>OK</message><skey>@crypt_3ea2fe08_723d1e1bd7b4171657b58c6d2849b367</skey><wxsid>9qxNHGgi9VP4/Tx6</wxsid><wxuin>979270107</wxuin><pass_ticket>%2BEdqKi12tfvM8ZZTdNeh4GLO9LFfwKLQRpqWk8LRYVWFkDE6%2FZJJXurz79ARX%2FIT</pass_ticket><isgrayscale>1</isgrayscale></error>
		*/
		var ret int = -1
		ret, this.state.wxSKey, this.state.wxsid, this.state.wxuin, this.state.wxPassTicket =
			parseTicket(resp.String())
		if ret != 0 {
			log.Println("failed")
		}
		this.reqState = REQ_WXINIT
		this.saveCookies(resp)
		this.eqch <- newEvent(EVT_REDIR_URL, []string{nsurl})
	}

}

func (this *longPoll) wxInit() {
	// # TODO: pass_ticket参数
	// nsurl = 'https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?r=1377482058764'
	// # v2 url:https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?r=-1222669677&lang=en_US&pass_ticket=%252BEdqKi12tfvM8ZZTdNeh4GLO9LFfwKLQRpqWk8LRYVWFkDE6%252FZJJXurz79ARX%252FIT
	// #nsurl = 'https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit?r=%s&lang=en_US&pass_ticket=' % \
	// #        (self.nowTime() - 3600 * 24 * 30)
	// #nsurl = self.urlBase + '/cgi-bin/mmwebwx-bin/webwxinit?r=%s&lang=en_US&pass_ticket=' % \
	// nsurl = self.urlBase + '/cgi-bin/mmwebwx-bin/webwxinit?r=%s&lang=en_US&pass_ticket=%s' % \
	// (self.nowTime() - 3600 * 24 * 30, self.wxPassTicket)
	// qDebug(nsurl)
	nsurl := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxinit?r=%d&lang=en_US&pass_ticket=%s",
		this.state.urlBase, time.Now().Unix()-3600*24*30, this.state.wxPassTicket)

	/*
		post_data = '{"BaseRequest":{"Uin":"%s","Sid":"%s","Skey":"","DeviceID":"%s"}}' % \
		(self.wxuin, self.wxsid, self.devid)

		req = requests.Request('post', nsurl, data=post_data.encode())
	*/

	postData := fmt.Sprintf(`{"BaseRequest":{"Uin":"%s","Sid":"%s","Skey":"","DeviceID":"%s"}}`,
		this.state.wxuin, this.state.wxsid, this.state.wxdevid)
	this.rops.JSON = postData

	resp, err := this.rses.Post(nsurl, this.rops)
	this.rops.JSON = nil
	if err != nil {
		log.Println(err, nsurl)
	}
	log.Println(resp.StatusCode, resp.Header, resp.Ok)
	bcc := resp.Bytes()
	this.saveContent("wxinit.json", bcc, resp, nsurl)
	defer resp.Close()

	if !resp.Ok {
		this.reqState = REQ_END
	} else {
		jcc, err := simplejson.NewJson(bcc)
		if err != nil {
			log.Println(err)
		} else {
			ret := jcc.GetPath("BaseResponse", "Ret").MustInt()
			log.Println("ret", ret)
			switch ret {
			case 1101:
				this.reqState = REQ_END
			default:
				this.state.wxSyncKey = jcc.Get("SyncKey")
				this.state.wxSKeyOld = this.state.wxSKey
				this.state.wxSKey = jcc.Get("SKey").MustString()
				if this.state.wxSKey != this.state.wxSKeyOld {
					log.Println("SKey updated:", this.state.wxSKeyOld, this.state.wxSKey)
				}
				this.state.wxInitRawData = resp.String()
				this.reqState = REQ_SYNC_CHECK
				this.reqState = REQ_CONTACT
				this.saveCookies(resp)
				this.eqch <- newEvent(EVT_GOT_BASEINFO, []string{resp.String()})
			}
		}
	}

}

func (this *longPoll) getContact() {

	// nsurl = 'https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?r=1377482079876'
	// #nsurl = 'https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?r='
	// nsurl = self.urlBase + '/cgi-bin/mmwebwx-bin/webwxgetcontact?r='
	nsurl := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxgetcontact?r=", this.state.urlBase)

	/*
		post_data = '{}'
		req = requests.Request('post', nsurl, data=post_data.encode())
	*/

	postData := fmt.Sprintf(`{}`)
	this.rops.JSON = postData

	resp, err := this.rses.Post(nsurl, this.rops)
	this.rops.JSON = nil
	if err != nil {
		log.Println(err, nsurl)
	}
	log.Println(resp.StatusCode, resp.Header, resp.Ok, len(resp.Bytes()))
	bcc := resp.Bytes()
	this.saveContent("wxcontact.json", bcc, resp, nsurl)
	defer resp.Close()

	if !resp.Ok {
		this.reqState = REQ_END
	} else {

		jcc, err := simplejson.NewJson(bcc)
		if err != nil {
			log.Println(err)
		} else {
			ret := jcc.GetPath("BaseResponse", "Ret").MustInt()
			log.Println("ret", ret)
			switch ret {
			case 1101:
				this.reqState = REQ_END
			default:
				this.state.wxContactRawData = resp.String()
				this.reqState = REQ_SYNC_CHECK
				// this.saveCookies(resp)
				this.eqch <- newEvent(EVT_GOT_CONTACT, []string{resp.String()})
			}
		}
	}
}

func (this *longPoll) packSyncKey() string {
	/*
		### make syncKey: format: 1_124125|2_452346345|3_65476547|1000_5643635
		syncKey = []
		for k in self.wxSyncKey['List']:
		elem = '%s_%s' % (k['Key'], k['Val'])
		syncKey.append(elem)

		# |需要URL编码成%7C
		syncKey = '%7C'.join(syncKey)   # [] => str''
	*/

	count := this.state.wxSyncKey.Get("Count").MustInt()
	log.Println("count:", count)
	skarr := make([]string, 0)
	for idx := 0; idx < count; idx++ {
		key := this.state.wxSyncKey.Get("List").GetIndex(idx).Get("Key").MustInt()
		val := this.state.wxSyncKey.Get("List").GetIndex(idx).Get("Val").MustInt()
		skarr = append(skarr, fmt.Sprintf("%d_%d", key, val))
	}
	return strings.Join(skarr, "%7C")
}

func parsesynccheck(str string) (code int, selector int) {
	// window.synccheck={retcode:"1101",selector:"0"}
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

func (this *longPoll) dumpState() {

	log.Println("qruuid		 ", this.state.qruuid)
	log.Println("qrpic		 ", len(this.state.qrpic))
	log.Println("logined	 ", this.state.logined)
	log.Println("redirUrl	 ", this.state.redirUrl)
	log.Println("urlBase	 ", this.state.urlBase)
	log.Println("pushUrlBase ", this.state.pushUrlBase)
	log.Println("wxdevid	 ", this.state.wxdevid)
	log.Println("wxuin		 ", this.state.wxuin)
	log.Println("wxsid		 ", this.state.wxsid)
	log.Println("wxDataTicket", this.state.wxDataTicket)
	log.Println("wxSKeyOld	 ", this.state.wxSKeyOld)
	log.Println("wxPassTicket", this.state.wxPassTicket)
	log.Println("wxuvid		 ", this.state.wxuvid)
	log.Println("wxAuthTicket", this.state.wxAuthTicket)
	log.Println("wxSyncKey	 ", this.state.wxSyncKey)
	log.Println("wxSKey		 ", this.state.wxSKey)
	log.Println("wxInitRawData  ", len(this.state.wxInitRawData))
	log.Println("wxContactRawData", len(this.state.wxContactRawData))

}

func (this *longPoll) syncCheck() {
	syncKey := this.packSyncKey()
	skey := strings.Replace(this.state.wxSKey, "@", "%40", -1)
	log.Println(this.state.wxSKey, "=>", skey)
	pass_ticket := strings.Replace(this.state.wxPassTicket, "%", "%25", -1)
	nsurl := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/synccheck?r=%d&skey=%s&sid=%s&uin=%s&deviceid=%s&synckey=%s&lang=en_US&pass_ticket=%s",
		this.state.pushUrlBase, this.nowTime(), skey, this.state.wxsid, this.state.wxuin,
		this.state.wxdevid, syncKey, pass_ticket)
	log.Println("requesting...", nsurl)

	resp, err := this.rses.Get(nsurl, this.rops)
	if err != nil {
		log.Println(err, nsurl)
	}
	log.Println(resp.StatusCode, resp.Header, resp.Ok)
	bcc := resp.Bytes()
	this.saveContent("synccheck.json", bcc, resp, nsurl)
	defer resp.Close()

	if !resp.Ok {
		this.reqState = REQ_END
	} else {
		log.Println(resp.String())
		retcode, selector := parsesynccheck(resp.String())

		switch retcode {
		case -1:
			log.Fatalln("wtf")
		case 1100:
			log.Println("maybe need reget SyncKey, rerun wxinit() ...")
		case 1101:
			this.dumpState()
			log.Println("maybe need rerun relogin...", resp.String())
			this.reqState = REQ_END
		case 0:
			switch selector {
			case 0: // go on syncCheck
			case 1:
				fallthrough
			case 2:
				fallthrough
			case 4:
				fallthrough
			case 5:
				fallthrough
			case 6:
				fallthrough
			case 7:
				this.reqState = REQ_WEB_SYNC
			default:
				log.Println("unknown selector:", retcode, selector)
			}
		default:
			log.Println("error sync check ret code:", retcode, selector)
		}
	}
}

func (this *longPoll) webSync() {
	skey := strings.Replace(this.state.wxSKey, "@", "%40", -1)
	log.Println(this.state.wxSKey, "=>", skey)
	pass_ticket := strings.Replace(this.state.wxPassTicket, "%", "%25", -1)
	nsurl := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxsync?sid=%s&skey=%s&lang=en_US&pass_ticket=%s", this.state.urlBase, this.state.wxsid, skey, pass_ticket)
	BaseRequest := map[string]string{
		"Uin":      this.state.wxuin,
		"Sid":      this.state.wxsid,
		"SKey":     this.state.wxSKey,
		"DeviceID": this.state.wxdevid}
	post_data_obj := simplejson.New()
	post_data_obj.Set("BaseRequest", BaseRequest)
	post_data_obj.Set("SyncKey", this.state.wxSyncKey)
	post_data_obj.Set("rr", this.nowTime())

	post_data_bin, err := post_data_obj.Encode()
	if err != nil {
		log.Println(err)
	}
	post_data := string(post_data_bin)

	this.rops.JSON = post_data
	this.rops.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	resp, err := this.rses.Post(nsurl, this.rops)
	this.rops.JSON = nil
	delete(this.rops.Headers, "Content-Type")
	if err != nil {
		log.Println(err, nsurl)
	}
	log.Println(resp.StatusCode, resp.Header, resp.Ok, len(resp.Bytes()))
	bcc := resp.Bytes()
	this.saveContent("websync.json", bcc, resp, nsurl)
	defer resp.Close()

	if !resp.Ok {
		this.reqState = REQ_END
	} else {
		jcc, err := simplejson.NewJson(bcc)
		if err != nil {
			log.Println(jcc)
			this.reqState = REQ_END
		} else {
			if jcc.GetPath("SyncKey", "Count").MustInt() == 0 {
				log.Println("websync's SyncKey empty, maybe need refresh...")
				this.reqState = REQ_END
			} else {
				// update SyncKey and SKey
				this.state.wxSyncKey = jcc.Get("SyncKey")

				// check data
				ret := jcc.GetPath("BaseResponse", "Ret").MustInt()
				switch ret {
				case 0:
					this.reqState = REQ_SYNC_CHECK
					this.saveCookies(resp)
				case 1101:
					log.Println("maybe need rerun refresh()...1101")
					this.reqState = REQ_END
				case -1:
					log.Println("wtf")
					this.reqState = REQ_END
				default:
					log.Println("web sync error:", ret)
					this.reqState = REQ_END
				}
			}
		}

		this.eqch <- newEvent(EVT_GOT_MESSAGE, []string{resp.String()})
	}

}
