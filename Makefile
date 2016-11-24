all:
	go install -v go-purple/purple

tox:
	go build -v -buildmode=c-shared -o libtox.so  ./tox-prpl

toxb:
	go build -v -o toxb.so ./tox-prpl

toxi:
	sudo install -v ./libtox.so /usr/lib/purple-2/

pld:
	go build -v -o prpld go-purple/purpled

wxp:
	go build -v -buildmode=c-shared -o libwechat.so  ./wechat-prpl

wxb:
	go build -v -o wechatb.so ./wechat-prpl

wxi:
	sudo install -v ./libwechat.so /usr/lib/purple-2/

wxl:
	go build -v ./wechat-prpl/wechat
	go install -v ./wechat-prpl/wechat

bot:
	go build -v -o yobot.bin go-purple/yobot/
