package session

import (
	"errors"
	"sync"
)

//MemorySession 内存Session
type MemorySession struct {
	sessionID string
	data      map[string]interface{}
	rwlock    sync.RWMutex
}

//NewMemorySession 构造函数
func NewMemorySession(ID string) *MemorySession {
	return &MemorySession{
		sessionID: ID,
		data:      make(map[string]interface{}, 16),
	}
}

func (m *MemorySession) Set(key string, value interface{}) (err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	m.data[key] = value
	return
}
func (m *MemorySession) Get(key string) (value interface{}, err error) {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()
	value, ok := m.data[key]
	if !ok {
		err = errors.New("memory session no this key ")
		return
	}
	return
}
func (m *MemorySession) Del(key string) (err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	delete(m.data, key)
	return
}
func (m *MemorySession) Save() (err error) {
	return
}
