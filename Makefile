all:
	go install -v
	go build -v
	go build -v examples/bot.go
