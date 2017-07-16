package fetchtitle

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"

	"github.com/PuerkitoBio/goquery"
	"github.com/kitech/colog"
	"gopkg.in/iconv.v1"
	// iconv2 "github.com/djimenez/iconv-go"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.76 Safari/537.36 Vivaldi/1.6.689.13"
const userAgent2 = "curl/7.54.0"

var userAgents = []string{userAgent, userAgent2}

func getUA() string { return userAgents[rand.Int()%len(userAgents)] }

var pxyurl string

func SetProxy(p string) {
	_, err := url.Parse(p)
	if err != nil {
		log.Println(err)
	} else {
		pxyurl = p
	}
}

// timeout seconds
// title, mime, error
func FetchMeta(u string, timeout int) (title string, mime string, err error) {
	title, mime, err = fetchWebMeta(u, timeout)
	if err != nil {
		rawerr := err
		errmsgs := []string{
			"connection refused", "Timeout",
			"connection reset by peer", "no such host",
		}
		for _, em := range errmsgs {
			if strings.Contains(err.Error(), em) {
				err = errors.New(fmt.Sprintf("Error: %s", em))
				break
			}
		}
		if false {
			log.Println(err.Error(), "<=", rawerr.Error())
		}
		return
	}

	// trim, join lines
	title = strings.Trim(title, "\n ")
	title = strings.Replace(title, "\n", "", 0)

	return
}

func fetchWebMeta(u string, timeout int) (title string, mime string, err error) {
	var cli *http.Client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		/*
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil
				// return http.ErrUseLastResponse
			},
		*/
	}
	var useProxy = len(pxyurl) > 0
	if useProxy {
		proxyUrl, err := url.Parse("http://127.0.0.1:8117")
		proxyUrl, err = url.Parse(pxyurl)
		if err != nil {
			log.Println("wtf", err)
		} else {
			cli.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyUrl)
		}
	}

	cli.Timeout = time.Duration(timeout) * time.Second

	var resp *http.Response
	for idx := 0; idx < 5; idx++ {
		var req *http.Request
		req, err = http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			// log.Println(err)
			return
		}

		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("User-Agent", getUA())
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Connection", "close")
		resp, err = cli.Do(req)
		if err != nil {
			// log.Println(err)
			return
		}

		// log.Println(resp.Request.URL.String(), req.URL.String())
		if resp.Request.URL.String() != u {
			// maybe redirected
			u = resp.Request.URL.String()
			resp.Body.Close()
			resp = nil
		} else {
			break
		}
	}
	if resp == nil {
		err = errors.New("reloc exceeded 5 times")
		return
	}
	defer resp.Body.Close()

	if false {
		bcc, _ := ioutil.ReadAll(resp.Body)
		log.Println(resp, string(bcc))
	}

	switch resp.StatusCode {
	case 300, 301, 302, 404, 406:
		title = fmt.Sprintf("%s", resp.Status)
		return
	}

	mime = resp.Header.Get("Content-Type")
	// log.Println(mime)

	// some other post reformat
	// image/video/audio/stream
	if len(mime) > 0 && !strings.Contains(mime, "text/") {
		sz := resp.ContentLength
		title = fmt.Sprintf("Size: %s", sizeToHuman(int64(sz)))
		return
	}
	defer resp.Body.Close()
	bcc, _ := ioutil.ReadAll(resp.Body)

	_, xcharset, _ := getCharset(bcc, mime)

	var hmime string
	var hcharset string
	title, hmime, hcharset, err = parseTitle(resp, bcc)
	if err != nil {
		// log.Println(err)
		return
	}
	if len(hmime) > len(mime) {
		mime = hmime
	}
	if hcharset == "" {
		hcharset = strings.ToLower(xcharset)
	}

	reparsetitle := false
	rebcc := []byte{}

	iconvif := func(cs string) bool {
		cseq1 := fmt.Sprintf("charset=%s", cs)
		cseq2 := fmt.Sprintf("; charset=\"%s\"", cs)

		if strings.Contains(strings.ToLower(mime), cseq1) ||
			strings.ToLower(hcharset) == cs ||
			bytes.Contains(bytes.ToLower(bcc), []byte(cseq2)) {
			ch, _ := iconv.Open("utf-8", cs)
			ntitle := ch.ConvString(title)
			if len(ntitle) > 0 {
				title = ntitle
			}
			ch.Close()
			reparsetitle = true
			ch, _ = iconv.Open("utf-8", cs)
			rebcc = make([]byte, len(bcc)*3)
			ch.Conv(bcc, rebcc)
			ch.Close()
			return true
		}
		return false
	}

	if iconvif("gbk") {
	} else if iconvif("gb2312") {
	} else {
	}

	log.Println(reparsetitle)
	// goquery需要utf-8的
	if reparsetitle {
		ntitle, nmime, nhcharset, nerr := parseTitle(resp, rebcc)
		if nerr == nil && title != ntitle {
			title, mime, hcharset = ntitle, nmime, nhcharset
		}
	}
	return
}

func parseTitle(resp *http.Response, bcc []byte) (
	title string, hmime string, hcharset string, err error) {

	reader := bytes.NewBuffer(bcc)
	doc, err := goquery.NewDocumentFromReader(reader)
	// doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		// log.Println(err)
		return "", "", "", err
	}

	sel := doc.Find("title")
	title, err = sel.Html()
	// log.Println(sel, title, err, sel.Length())
	if sel.Length() >= 0 {
		headhtml, err := doc.Find("head").Html()
		headhtml, err = doc.Html()
		log.Println(doc.Length(), err, headhtml)
		log.Println(resp.StatusCode, resp.Status, resp.Header)
		// ioutil.WriteFile("hehhe.html", []byte(headhtml), 0755)
	}

	sel = doc.Find("meta")
	sel.Each(func(idx int, s *goquery.Selection) {
		if _, ok := s.Attr("http-equiv"); ok {
			if mime, ok := s.Attr("content"); ok {
				hmime = mime
			}
		}
		if cs, ok := s.Attr("charset"); ok {
			hcharset = cs
		}
	})

	return title, hmime, hcharset, nil
}

func getCharset(content []byte, contentType string) (
	e encoding.Encoding, name string, certain bool) {
	e, name, certain = charset.DetermineEncoding(content, contentType)
	log.Println(e, name, certain)
	return
}

func sizeToHuman(sz int64) (hsz string) {
	if sz > 1024*1024*1024 {
		hsz = fmt.Sprintf("%.2fG", float64(sz)/(1024*1024*1024.0))
	} else if sz > 1024*1024 {
		hsz = fmt.Sprintf("%.1fM", float64(sz)/(1024*1024.0))
	} else if sz > 1024 {
		hsz = fmt.Sprintf("%.1fK", float64(sz)/(1024.0))
	} else {
		hsz = fmt.Sprintf("%dB", sz)
	}
	return
}

func init() {
	colog.Register()
	colog.SetFlags(log.Flags() | log.Lshortfile | log.LstdFlags)
	time.Sleep(0)
}
