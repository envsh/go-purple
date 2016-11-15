all:
	go install -v yobot/purple
	go build -v

tox:
	go build -v -buildmode=c-shared -o libtox.so  ./tox-prpl

toxb:
	go build -v -o toxb.so ./tox-prpl

toxi:
	sudo install -v ./libtox.so /usr/lib/purple-2/libtox.so

pld:
	go build -v -o prpld yobot/purpled

