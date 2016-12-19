package fetchtitle

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kitech/colog"
	"gopkg.in/iconv.v1"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.76 Safari/537.36 Vivaldi/1.6.689.13"

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
		req.Header.Set("Accept", "*.*")
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
	case 300, 301, 302, 404:
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

	var hmime string
	title, hmime, err = parseTitle(resp)
	if err != nil {
		// log.Println(err)
		return
	}
	if len(hmime) > len(mime) {
		mime = hmime
	}

	// TODO 自动检测编码并尝试转码
	if strings.Contains(strings.ToLower(mime), "charset=gbk") {
		ch, _ := iconv.Open("utf-8", "gbk")
		ntitle := ch.ConvString(title)
		if len(ntitle) > 0 {
			title = ntitle
		}
		ch.Close()
	} else if strings.Contains(strings.ToLower(mime), "charset=gb2312") {
		ch, _ := iconv.Open("utf-8", "gb2312")
		ntitle := ch.ConvString(title)
		if len(ntitle) > 0 {
			title = ntitle
		}
		ch.Close()
	}

	return
}

func parseTitle(resp *http.Response) (string, string, error) {
	if false {
		buf := make([]byte, 8192)
		n, err := resp.Body.Read(buf)
		if err != nil {
			log.Println(err, n)
		}
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		// log.Println(err)
		return "", "", err
	}

	sel := doc.Find("title")
	title, err := sel.Html()
	// log.Println(sel, title, err, sel.Length())

	hmime := ""
	sel = doc.Find("meta")
	sel.Each(func(idx int, s *goquery.Selection) {
		if _, ok := s.Attr("http-equiv"); ok {
			if mime, ok := sel.Attr("content"); ok {
				hmime = mime
			}
		}
	})

	log.Println(resp.Close)
	resp.Body.Close()
	resp = nil
	return title, hmime, nil
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
