package main

type Backend interface {
	isconnected() bool
	disconnect()
}

type BackendBase struct {
	ctx    *Context
	conque chan interface{}
	proto  string
	name   string
}
