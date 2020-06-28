package session

//SessionMgr 定义管理者  管理所有的session
type SessionMgr interface {
	//初始化
	Init(addr string, option ...string) (err error)
	CreatSession() (session Session, err error)
	GetSession(sessionID string) (session Session, err error)
}
