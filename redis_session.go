package session

import (
	"encoding/json"
	"sync"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

//使用常量来定义内存中的map是否被操作
const (
	SessionFlagNone = iota
	SessionFlagModify
)

// RedisSession 结构图
type RedisSession struct {
	sessionID string
	//设置session，可以先放在内存的map中；再批量导入redis中，提升性能
	sessionMap map[string]interface{}
	rwlock     sync.RWMutex
	pool       *redis.Pool
	flag       int //记录内存中的map是否被操作
}

//NewRedisSession  构造函数 返回RedisSession结构体
func NewRedisSession(id string, pool *redis.Pool) *RedisSession {
	return &RedisSession{
		sessionID:  id,
		pool:       pool,
		flag:       SessionFlagNone,
		sessionMap: make(map[string]interface{}, 16),
	}
}

//Set session存储 到内存中的map
func (r *RedisSession) Set(key string, value interface{}) (err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	r.sessionMap[key] = value
	r.flag = SessionFlagModify
	return
}

//Save session存储 从内存中存到redis中
func (r *RedisSession) Save() (err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	if r.flag != SessionFlagModify { //数据没变  不需要存
		return
	}
	byteData, err := json.Marshal(r.sessionMap)
	if err != nil {
		return
	}
	//获取redis连接
	conn := r.pool.Get()
	_, err = conn.Do("SET", r.sessionID, string(byteData))
	r.flag = SessionFlagNone //修改状态
	if err != nil {
		return
	}
	return
}

//Get 根据key 获取map数据中的value(先从redis中加载到内存里面)
func (r *RedisSession) Get(key string) (value interface{}, err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	//先判断内存中是否有数据
	err = r.loadFromRedis()
	if err != nil {
		return
	}
	value, ok := r.sessionMap[key]
	if !ok {
		err = errors.New("key not exit")
		return
	}
	return
}


//从redis中加载数据
func (r *RedisSession) loadFromRedis() (err error) {
	conn := r.pool.Get()
	reply, err := conn.Do("GET", r.sessionID)
	if err != nil {
		return
	}
	//转字符串
	data, err := redis.String(reply, err)
	if err != nil {
		return
	}
	//取到的东西反序列化到内存中的map
	err = json.Unmarshal([]byte(data), &r.sessionMap)
	if err != nil {
		return
	}
	return
}

//Del 删除
func (r *RedisSession) Del(key string) (err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	r.flag = SessionFlagModify
	delete(r.sessionMap, key)
	return
}
