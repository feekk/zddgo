package zddgo

import(
	"sync"
	"github.com/feekk/zddgo/redis"
)

var(
	defaultConn *redis.RedisPool 
	otherConns *sync.Map = &sync.Map{}
)

func RedisInit(c *RedisConf) (err error){
	if !c.Use {
		return
	}
	if defaultConn, err = NewRedisPool(&c.Default); err != nil {
		return
	}
	
	if len(c.Connection) <= 0 {
		return
	}
	var conn *redis.RedisPool
	for _, conConf := range c.Connection{
		conn, err = NewRedisPool(&conConf)
		if err != nil{
			return 
		}
		otherConns.Store(conConf.Name, *conn)
	} 
	return
}
//
// conn := cache.Def().Get()
// default conn.Close()
//
func RedisDef() *redis.RedisPool{
	return defaultConn
}
//
// conn := cache.Conn("name1").Get()
// default conn.Close()
//
func RedisConn(key string) (*redis.RedisPool, bool){
	if pool, ok := otherConns.Load(key); ok {
		 p := pool.(redis.RedisPool)
		return &p, true
	}
	return nil, false
}