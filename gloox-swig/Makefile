all:
    # -small
	swig -c++ -go -cgo -intgosize 64 -v  -Wall -I/usr/include -I/usr/include/gloox gloox.swig_
	# sed -i -E  "s/arg([0-9]+) = \(va_list\)\*argp([0-9]+)/memcpy\(\&arg\1, argp\2, sizeof\(va_list)\\)/" gloox_wrap.c
	# sed -i -E 's/func Purple_([a-z])/func \U\1/' gloox.go
	# sed -i -E 's/gloox::XMLNS_JABBER_RPC/gloox::XMLNS_NICKNAME/' gloox_wrap.cxx
	go build -v -x

build:
	go build -v -x
b:
	go build -v -x

install:
	go install -v
i:
	go install -v

test:
	go test -v -run 1

