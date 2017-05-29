package fetchtitle

import (
	"fmt"
	"log"
	"testing"
)

func TestFetchTitle(t *testing.T) {
	timeout := 10

	urls := []string{
		"http://lilydjwg.is-programmer.com/",
		"http://www.baidu.com",
		"https://zh.wikipedia.org", // redirection
		"http://redis.io/",
		"http://lilydjwg.is-programmer.com/2012/10/27/streaming-gzip-decompression-in-python.36130.html",                // maybe timeout
		"http://img.vim-cn.com/22/cd42b4c776c588b6e69051a22e42dabf28f436",                                               // image with length
		"https://github.com/m13253/titlebot/blob/master/titlebot.py_",                                                   // 404
		"http://lilydjwg.is-programmer.com/admin",                                                                       // redirection
		"http://twitter.com",                                                                                            // connect timeout
		"http://www.wordpress.com",                                                                                      // reset
		"http://jquery-api-zh-cn.googlecode.com/svn/trunk/xml/jqueryapi.xml",                                            // xml
		"http://lilydjwg.is-programmer.com/user_files/lilydjwg/config/avatar.png",                                       // PNG
		"http://img01.taobaocdn.com/bao/uploaded/i1/110928240/T2okG7XaRbXXXXXXXX_!!110928240.jpg",                       // JPEG with Start Of Frame as the second block
		"http://file3.u148.net/2013/1/images/1357536246993.jpg",                                                         // JPEG that failed previous code
		"http://gouwu.hao123.com/",                                                                                      // HTML5 GBK encoding
		"https://github.com/lilydjwg/winterpy",                                                                          // github url finder
		"http://github.com/lilydjwg/winterpy",                                                                           // github url finder with redirect
		"http://导航.中国/",                                                                                                 // Punycode. This should not be redirected
		"http://t.cn/zTOgr1n",                                                                                           // multiple redirections
		"http://www.galago-project.org/specs/notification/0.9/x408.html",                                                // </TITLE\n>
		"http://x.co/dreamz",                                                                                            // redirection caused false ConnectionClosed error
		"http://www.83wyt.com",                                                                                          // reversed meta attribute order
		"https://www.inoreader.com",                                                                                     // malformed start tag: <meta http-equiv="Content-Type" content="text/html" ; charset="UTF-8">
		"https://linuxtoy.org/archives/linux-deepin-2014-alpha-into-new-deepin-world.html",                              // charref outside ASCII
		"http://74.125.235.191/search?site=&source=hp&q=%E6%9C%8D%E5%8A%A1%E5%99%A8+SSD&btnG=Google+%E6%90%9C%E7%B4%A2", // right charset in HTTP, wrong in HTML
		"http://digital.sina.com.hk/news/-7-1514837/1.html",                                                             // mixed Big5 and non-Big5 escaped Unicode character
		"http://cache.baiducontent.com/c?m=9f65cb4a8c8507ed4fece7631046893b4c4380147c808c5528888448e435061e5a27b9e867750d04d6c57f6102ad4b57f7fa3372340126bc9fcc825e98e6d27e20d77465671df65663a70edecb5124b137e65ffed86ef0bb8025e3ddc5a2de4352ba44757d97818d4d0164dd1efa034093b1e842022e60adec40728f2d6058e93430c6508ae5256f779686d94b3db3&p=882a9e41c0d25ffc57efdc394c52&newp=8a64865b85cc43ff57e6902c495f92695803ed603fd3d7&user=baidu&fm=sc&query=mac%CF%C2%D7%EE%BA%C3%B5%C4%C8%CB%C8%CB%BF%CD%BB%A7%B6%CB&qid=&p1=5", // HTML document inside another, correct charset is in outside one and title inside
		"http://www.wooyun.org/bugs/wooyun-2014-069132?214&1780",                 // provide cookie or loop
		"https://togetherjs.com/",                                                // fail if no SNI
		"https://forum.suse.org.cn/",                                             // most Linux distributions seem not to trust this
		"http://www.aosabook.org/en/posa/parsing-xml-at-the-speed-of-light.html", // <span> tag inside <title> tag
		"http://github.com/contact",                                              // redirect and should not use GitHub API
		"http://caml.inria.fr/pub/docs/manual-ocaml/",
		"http://localhost/", // should fail with ValueError
	}

	urls = []string{
		// "https://gist.github.com/4caa949091e80834c965156f0e3b493c", //empty title?
		// "https://dn-fishroom.qbox.me/sticker/1370.png",
		// "http://fishroom.tuna.moe/log/linuxba/2016-12-18/553",
		// "https://github.com/kitech/go-toxcore",
		"http://wuxia.qq.com/main.shtml",
	}

	titles := []string{}
	// SetProxy("http://127.0.0.1:8117")

	for idx, url := range urls {
		t.Run(fmt.Sprintf("url%02d", idx), func(t *testing.T) {
			title, mime, err := FetchMeta(url, timeout)
			if err != nil {
				log.Println(err, title, mime, url)
				t.Fail()
			} else {
				log.Printf("id:%d, mime:%s, title:%s, url: %s\n",
					idx, mime, title, url)
				log.Println("title len:", len(title), title)
				if idx < len(titles) && title != titles[idx] {
					log.Println(title, "!=", titles[idx], url)
					t.Fail()
				}
			}
		})
	}

}
