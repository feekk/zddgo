package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// Implement pool.Poolable & redis.Conn interface
type RedisPooledConnection struct {
	*PooledObject
	addr string
	c    redis.Conn
	dp   *Pool
}

type RedisPooledConnFactory struct {
	Password       string
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

var DefaultRedisPooledConnFactory RedisPooledConnFactory = RedisPooledConnFactory{
	ConnectTimeout: 50 * time.Millisecond,
	ReadTimeout:    100 * time.Millisecond,
	WriteTimeout:   100 * time.Millisecond,
}

// Factory to create new connection
func (f RedisPooledConnFactory) Create(address string) (pc Poolable, err error) {
	c, err := redis.DialTimeout("tcp", address, f.ConnectTimeout, f.ReadTimeout, f.WriteTimeout)
	if err != nil {
		return nil, err
	}

	if len(f.Password) > 0 {
		if _, err = c.Do("AUTH", f.Password); err != nil {
			c.Close()
			return nil, err
		}
	}

	pc = &RedisPooledConnection{
		PooledObject: &PooledObject{},
		c:            c,
		addr:         address,
	}
	return pc, nil
}

func (f RedisPooledConnFactory) Validate(pc Poolable) (err error) {
	_, err = pc.(*RedisPooledConnection).Do("PING")
	return err
}

func (f RedisPooledConnFactory) Close(pc Poolable) error {
	return pc.(*RedisPooledConnection).c.Close()
}

// Close closes the connection.
// @Override
func (pc *RedisPooledConnection) Close() error {
	err := pc.c.Err()
	if pc.dp == nil {
		return pc.c.Close()
	}
	if err != nil && err != redis.ErrNil {
		return pc.dp.Put(pc, true)
	} else {
		return pc.dp.Put(pc, false)
	}
}

// Err returns a non-nil value if the connection is broken. The returned
// value is either the first non-nil value returned from the underlying
// network connection or a protocol parsing error. Applications should
// close broken connections.
// @Override
func (pc *RedisPooledConnection) Err() error {
	return pc.c.Err()
}

// Do sends a command to the server and returns the received reply.
// @Override
func (pc *RedisPooledConnection) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return pc.c.Do(commandName, args...)
}

// Send writes the command to the client's output buffer.
// @Override
func (pc *RedisPooledConnection) Send(commandName string, args ...interface{}) error {
	return pc.c.Send(commandName, args...)
}

// Flush flushes the output buffer to the Redis server.
// @Override
func (pc *RedisPooledConnection) Flush() error {
	return pc.c.Flush()
}

// Receive receives a single reply from the Redis server
// @Override
func (pc *RedisPooledConnection) Receive() (reply interface{}, err error) {
	return pc.c.Receive()
}

// 用于兼容redigo接口拷贝自redigo代码
type errorConnection struct{ err error }

func (ec errorConnection) Do(string, ...interface{}) (interface{}, error) { return nil, ec.err }
func (ec errorConnection) Send(string, ...interface{}) error              { return ec.err }
func (ec errorConnection) Err() error                                     { return ec.err }
func (ec errorConnection) Close() error                                   { return ec.err }
func (ec errorConnection) Flush() error                                   { return ec.err }
func (ec errorConnection) Receive() (interface{}, error)                  { return nil, ec.err }

type RedisPool struct {
	dp *Pool
}

func NewRedisPool(servers []string, connFactory PooledConnFactory, poolConfig PoolConfig) *RedisPool {
	return &RedisPool{dp: NewPool(servers, connFactory, poolConfig)}
}

// 兼容redigo pool的Get()接口
func (p *RedisPool) Get() redis.Conn {
	raw, err := p.dp.Get()
	if err != nil {
		return errorConnection{err}
	}

	pc := raw.(*RedisPooledConnection)
	pc.dp = p.dp
	return pc
}

// 兼容redigo的关闭函数接口
func (p *RedisPool) Close() (err error) {
	p.dp.Shutdown()
	return
}

func (p *RedisPool) GetPoolStats() (stats []PoolStats) {
	return p.dp.GetPoolStats()
}
