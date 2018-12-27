package main

import (
	"net"
	"sync"
)

type ConnPool struct {
	conns  map[string][]net.Conn
	locker *sync.RWMutex
}

func NewConnPool() (cp *ConnPool) {
	cp = &ConnPool{}
	cp.conns = map[string][]net.Conn{}
	cp.locker = &sync.RWMutex{}
	return
}

func (self *ConnPool) GetConn(key string) (conn net.Conn) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if cs, ok := self.conns[key]; ok {
		conn = cs[0]
		cs = cs[1:]

		if len(cs) <= 0 {
			delete(self.conns, key)
		} else {
			self.conns[key] = cs
		}
	}
	return
}

func (self *ConnPool) AddConn(key string, conn net.Conn) {
	self.locker.Lock()
	defer self.locker.Unlock()

	cs := self.conns[key]
	if cs == nil {
		cs = []net.Conn{}
	}
	cs = append(cs, conn)

	self.conns[key] = cs
	return
}

func (self *ConnPool) AddConns(key string, conns []net.Conn) {
	self.locker.Lock()
	defer self.locker.Unlock()

	cs := self.conns[key]
	if cs == nil {
		cs = []net.Conn{}
	}
	cs = append(cs, conns...)

	self.conns[key] = cs
	return
}
