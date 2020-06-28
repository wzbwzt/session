package session

import "fmt"

//中间件让用户选择使用哪个版本
var (
	sessionmgr SessionMgr
)

func Init(provider, addr string, option ...string) (err error) {
	switch provider {
	case "memary":
		sessionmgr = NewMemorySessionMgr()
	case "redis":
		sessionmgr = NewRedisSessionMgr()
	default:
		fmt.Errorf("不支持")
		return
	}
	err = sessionmgr.Init(addr, option...)
	return
}
