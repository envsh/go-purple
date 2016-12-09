package main

type Backend interface {
}

type BackendBase struct {
	ctx    *Context
	conque chan interface{}
	proto  int
	name   string
}
