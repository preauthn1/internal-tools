package util

import (
	"sync"
)

// SafeMap 安全的map對象
type SafeMap struct {
	sync.RWMutex
	innerMap map[interface{}]interface{}
}

// NewSafeMap 創建safemap
func NewSafeMap() *SafeMap {
	sm := new(SafeMap)
	sm.innerMap = make(map[interface{}]interface{})

	return sm
}

// Get ?
func (m *SafeMap) Get(key interface{}) interface{} {
	m.RLock()
	defer m.RUnlock()

	return m.innerMap[key]
}

// Set ?
func (m *SafeMap) Set(key interface{}, value interface{}) {
	m.Lock()
	defer m.Unlock()

	m.innerMap[key] = value
}

// Del ?
func (m *SafeMap) Del(key interface{}) {
	m.Lock()
	defer m.Unlock()

	delete(m.innerMap, key)
}

// Len ?
func (m *SafeMap) Len() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.innerMap)
}

// RLockRange ?
func (m *SafeMap) RLockRange(f func(interface{}, interface{})) {
	m.RLock()
	defer m.RUnlock()

	for k, v := range m.innerMap {
		f(k, v)
	}
}

// LockRange ?
func (m *SafeMap) LockRange(f func(interface{}, interface{})) {
	m.Lock()
	defer m.Unlock()

	for k, v := range m.innerMap {
		f(k, v)
	}
}
