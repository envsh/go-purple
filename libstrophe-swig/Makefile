all:
	swig -go -cgo -intgosize 64 -v  -Wall -I/usr/include/ libstrophe.swig_
	# sed -i -E  "s/arg([0-9]+) = \(va_list\)\*argp([0-9]+)/memcpy\(\&arg\1, argp\2, sizeof\(va_list)\\)/" purple_wrap.c
	# sed -i -E 's/func Purple_([a-z])/func \U\1/' purple.go
	go build -v -x

