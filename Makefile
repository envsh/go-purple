all:
	go install -v -x
	go build -v
	go build -v examples/bot.go
