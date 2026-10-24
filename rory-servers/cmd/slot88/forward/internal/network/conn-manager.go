package network

import (
	"sync"
)

var connm *connmanager

func init() {
	connm = new(connmanager)
	connm.m = make(map[string]*Conn)
}

type connmanager struct {
	sync.RWMutex
	m map[string]*Conn
}

func (m *connmanager) get(addr string) *Conn {
	m.RLock()
	defer m.RUnlock()

	return m.m[addr]
}

func (m *connmanager) set(addr string, conn *Conn) {
	m.Lock()
	defer m.Unlock()

	m.m[addr] = conn
}
