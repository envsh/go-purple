all:
	go install -v yobot/purple
	go build -v

tox:
	go build -v -buildmode=c-shared -o tox.so  ./tox-prpl
