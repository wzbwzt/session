package session

import (
	"errors"
	"sync"

	uuid "github.com/satori/go.uuid"
)

//MemorySessionMgr 结构体对象
type MemorySessionMgr struct {
	sessionMap map[string]Session
	rwlock     sync.RWMutex
}

//NewMemorySessionMgr  构造函数
func NewMemorySessionMgr() *MemorySessionMgr {
	return &MemorySessionMgr{
		sessionMap: make(map[string]Session, 1024),
	}
}

func (m *MemorySessionMgr) Init(addr string, option ...string) (err error) {
	return
}

func (m *MemorySessionMgr) CreatSession() (session Session, err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	id:= uuid.NewV4()
	sessionID := id.String()
	session = NewMemorySession(sessionID)
	m.sessionMap[sessionID] = session
	return
}

func (m *MemorySessionMgr) GetSession(sessionID string) (session Session, err error) {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()
	session, ok := m.sessionMap[sessionID]
	if !ok {
		err = errors.New("mgr has not this sessionID")
		return
	}
	return session,nil
}
