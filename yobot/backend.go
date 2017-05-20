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
	rname  string // maybe invalid name for some backend
	uid    string
}

func (this *BackendBase) fmtname() string {
	return this.name
}
