package zddgo

import(
	"time"
	"strings"
	"github.com/feekk/zddgo/redis"
	"github.com/feekk/zddgo/ztime"
	"github.com/feekk/zddgo/errors"
)


func NewRedisPool(c *RedisPoolConfig) (p *redis.RedisPool, err error){
	servers := strings.Split(c.Dsn, ";")
	factory := redis.RedisPooledConnFactory{
		ConnectTimeout: time.Duration(c.ConnectTimeout),
		ReadTimeout: time.Duration(c.ReadTimeout),
		WriteTimeout: time.Duration(c.WriteTimeout),
	}
	//check
	var pc redis.Poolable
	for _, server := range servers{
		pc, err = factory.Create(server)
		if err != nil{
			err = errors.With(err)
			return
		}
		if err = factory.Validate(pc); err != nil {
			err = errors.With(err)
			return
		}
		factory.Close(pc)
	}
	//instance
	conf := redis.PoolConfig{
		MaxIdle: c.MaxIdle,
		MaxActive: c.MaxActive,
		IdleTimeout: time.Duration(c.IdleTimeout),
		MaxFails:5,
	}
	if c.Pwd != "" {
		factory.Password = c.Pwd
	}

	p = redis.NewRedisPool(servers, factory, conf)
	return 
}


type RedisPoolConfig struct{
	Name string
	Dsn string
	Pwd string
	Db int
	MaxIdle int
	MaxActive int
	IdleTimeout ztime.Duration //idle timeout
	ConnectTimeout ztime.Duration
	ReadTimeout ztime.Duration //read timeout
	WriteTimeout ztime.Duration //write timeout
}