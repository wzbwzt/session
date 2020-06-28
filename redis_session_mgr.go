package session

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

//RedisSessionMgr 结构体对象
type RedisSessionMgr struct {
	addr       string
	passwd     string
	pool       *redis.Pool
	sessionMap map[string]Session
	rwlock     sync.RWMutex
}

//NewRedisSessionMgr 构造函数
func NewRedisSessionMgr() *RedisSessionMgr {
	return &RedisSessionMgr{
		sessionMap: make(map[string]Session, 32),
	}

}

//创建redis连接池
func myPool(addr, passwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     64,
		MaxActive:   100,
		IdleTimeout: 24 * time.Second,
		Dial:func()(redis.Conn,error){
			conn,err:=	redis.Dial("tcp",addr)
			if err != nil {
				return nil, err
			}
			//若有密码，判断
			if _,err=conn.Do("AUTH",passwd);err!=nil{
				conn.Close()
				return nil ,err
			}
			return conn,err
	},
	//连接测试，开发时写
	//上线后注释掉
	TestOnBorrow:func(conn redis.Conn,t time.Time) error{
		_,err:=conn.Do("PING")
		return err
	},
}
}
//Init 初始化 加载 redis 地址
func (r *RedisSessionMgr) Init(addr string, option ...string) (err error) {
	//若有其他option
	if len(option) > 0 {
		r.passwd = option[0]
	}
	//创建连接池
	r.pool = myPool(addr, r.passwd)
	r.addr = addr
	return
}

//CreatSession 创建一个session
func (r *RedisSessionMgr) CreatSession() (session Session, err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	id:= uuid.NewV4()
	sessionID:=id.String()
	session= NewRedisSession(sessionID,r.pool)
	r.sessionMap[sessionID]=session
	return

}

//GetSession 获取一个session
func (r *RedisSessionMgr) GetSession(sessionID string) (session Session, err error) {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()
	session, ok := r.sessionMap[sessionID]
	if !ok {
		err = errors.New("mgr has not this sessionID")
		return
	}
	return session,nil
}
