package redis

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// If an attempt is made to return an object to the pool that is in any state
// other than allocated (i.e. borrowed).
// Attempting to return an object more than once or attempting to return an
// object that was never borrowed from the pool will trigger this error.
var ErrReturnInvalid error = errors.New("pool: object has already been returned to this pool or is invalid")

type Pool struct {
	// Server address list, e.t. []string{"127.0.0.1:8080", "127.0.0.1:8081"}
	serverList []string
	// Function to create a new pooled client for server @serverAddr
	connFactory PooledConnFactory
	// Sharded pool by server address
	poolShards []*PoolShard
	// Pool configuration, e.t. maxIdle, maxActive, ...
	poolConfig PoolConfig
	// @atomic index to pick the next shard
	index uint32
	// Suspect shards, should be checked immediately
	suspectShards chan *PoolShard
	// Current available servers
	numAvailable int
	maxRetry int
	// Stopper signal to stop checker coroutine
	stopper chan struct{}
	// WaitGroup to wait health checker goroutine to stop
	wg sync.WaitGroup
}

func NewPool(servers []string, connFactory PooledConnFactory, poolConfig PoolConfig) *Pool {
	if servers == nil || len(servers) == 0 || connFactory == nil {
		panic("Illegal Arguments")
	}

	numServers := len(servers)

	dp := &Pool{
		serverList:    servers,
		connFactory:   connFactory,
		poolConfig:    poolConfig,
		suspectShards: make(chan *PoolShard, 100),
		numAvailable:  numServers,
		maxRetry:      numServers,
		stopper:       make(chan struct{}),
	}
	poolShards := make([]*PoolShard, numServers)
	for i := 0; i < numServers; i++ {
		shard := NewPoolShard(servers[i], dp, poolConfig.MaxIdle, poolConfig.MaxActive, poolConfig.MaxFails)
		poolShards[i] = shard
	}
	dp.poolShards = poolShards
	if numServers < 5 {
		dp.maxRetry = 5
	}

	go dp.goCheckServer()

	return dp
}

func (dp *Pool) Get() (Poolable, error) {
	var localIdx uint32 = atomic.AddUint32(&dp.index, 1)

	for tries := 0; tries < dp.maxRetry; tries++ {
		idx := (localIdx + uint32(tries)) % uint32(len(dp.serverList))

		// The server shard selected may be down, continue to get the next one
		if !dp.poolShards[idx].isAvailable() {
			atomic.AddUint32(&dp.index, 1)
			continue
		}

		c, err := dp.poolShards[idx].get()
		if err != nil {
			atomic.AddUint32(&dp.index, 1)
			continue
		}
		return c, nil
	}

	return nil, fmt.Errorf("pool: failed to get connection after %d retries", dp.maxRetry)
}

func (dp *Pool) Put(c Poolable, broken bool) error {
	shard := c.getDataSource()
	if shard == nil {
		dp.connFactory.Close(c)
		return ErrReturnInvalid
	}
	return shard.put(c, broken)
}

func (dp *Pool) markAvailable(shard *PoolShard, b bool) {
	if b {
		if shard.markAvailable(true) {
			dp.numAvailable++
		}
	} else {
		if !shard.isAvailable() {
			return
		}
		totalServers := len(dp.serverList)
		// Ensure that at most 1/3 servers can be marked as unavaialable
		if dp.numAvailable*3 > totalServers*2 {
			if shard.markAvailable(false) {
				dp.numAvailable--
			}
		} else {
		}
	}
}

func (dp *Pool) checkServer(server string) (ok bool) {
	for tries := 1; tries <= 2; tries++ {
		c, err := dp.connFactory.Create(server)
		if err != nil {
			continue
		}

		if err := dp.connFactory.Validate(c); err != nil {
			dp.connFactory.Close(c)
			continue
		}

		dp.connFactory.Close(c)
		return true
	}

	return false
}

// Check server availiability periodically
func (dp *Pool) goCheckServer() {
	defer dp.wg.Done()
	dp.wg.Add(1) // XXX
	var timer *time.Ticker = time.NewTicker(3 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			for _, shard := range dp.poolShards {
				// Healthy shards can be exampt from examination
				if !shard.suspectable() && shard.isAvailable() {
					continue
				}
				dp.markAvailable(shard, dp.checkServer(shard.server))
			}
		case shard := <-dp.suspectShards:
			dp.markAvailable(shard, dp.checkServer(shard.server))
		case <-dp.stopper:
			return
		}
	}
}

func (dp *Pool) Shutdown() {
	close(dp.stopper)
	dp.wg.Wait()
	for _, shard := range dp.poolShards {
		shard.Close()
	}
}

type PoolStats struct {
	Shard        string `json:"shard"`
	Available    bool   `json:"available"`
	NumActive    int    `json:"num_active"`
	NumGet       uint64 `json:"num_get"`
	NumPut       uint64 `json:"num_put"`
	NumBroken    uint64 `json:"num_broken"`
	NumDial      uint64 `json:"num_dial"`
	NumDialError uint64 `json:"num_dial_error"`
	NumEvict     uint64 `json:"num_evict"`
	NumClose     uint64 `json:"num_close"`
}

func (dp *Pool) GetPoolStats() (stats []PoolStats) {
	stats = make([]PoolStats, len(dp.serverList))
	for i, shard := range dp.poolShards {
		stats[i] = shard.getStats()
	}
	return
}


type PoolConfig struct {
	// Maximum idle connections per shard
	MaxIdle int
	// The maximum number of active connections that can be allocated from per pool shard at the same time.
	// The default value is 100
	MaxActive int
	// Timeout to evict idle connections
	IdleTimeout time.Duration
	// Test if connection broken on borrow
	// If set this flag, the "test" function should also provided.
	TestOnBorrow bool
	// Number of max fails threshold to triger health check
	MaxFails int
}

var DefaultPoolConfig PoolConfig = PoolConfig{
	MaxIdle:     50,
	MaxActive:   100,
	IdleTimeout: 300 * time.Second,
	MaxFails:    5,
}



type PooledConnFactory interface {
	// Function to create a new pooled client for server @serverAddr
	Create(addr string) (Poolable, error)

	// testOnBorrow is an optional application supplied function for checking
	// the health of an idle connection before the connection is used again by
	// the application. Argument t is the time that the connection was returned
	// to the pool. If the function returns an error, then the connection is
	// closed.
	Validate(c Poolable) error

	// Function to destroy a connection
	Close(c Poolable) error
}


// Poolable represents a connection to a server.
type Poolable interface {
	// @private
	lock()
	unlock()

	// @private
	setTime(t time.Time)
	getTime() time.Time

	// @private
	getDataSource() *PoolShard
	setDataSource(shard *PoolShard)

	// @private
	isBorrowed() bool
	setBorrowed(b bool)
}

// Implements Poolable interface
type PooledObject struct {
	dataSource *PoolShard
	t          time.Time
	borrowed   bool
	state      int
	mu         sync.Mutex
}

// @private
func (pc *PooledObject) lock() {
	pc.mu.Lock()
}

// @private
func (pc *PooledObject) unlock() {
	pc.mu.Unlock()
}

// @private
func (pc *PooledObject) setDataSource(shard *PoolShard) {
	pc.dataSource = shard
}

// @private
func (pc *PooledObject) getDataSource() (shard *PoolShard) {
	return pc.dataSource
}

// @private
func (pc *PooledObject) setTime(t time.Time) {
	pc.t = t
}

// @private
func (pc *PooledObject) getTime() (t time.Time) {
	t = pc.t
	return t
}

// @private
func (pc *PooledObject) isBorrowed() bool {
	return pc.borrowed
}

// @private
func (pc *PooledObject) setBorrowed(b bool) {
	pc.borrowed = b
}








var (
	// errPoolExhausted is returned from a pool connection method when
	// the maximum number of connections in the pool has been reached.
	ErrPoolExhausted = errors.New("dpool: connection pool exhausted")

	ErrPoolClosed = errors.New("dpool: connection pool closed")
	ErrConnClosed = errors.New("dpool: connection closed")
)

type PoolShard struct {
	// Maximum number of idle connections in the pool.
	// @const
	maxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	// @const
	maxActive int32

	// Current number of active connections
	// @atomic
	active int32

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	idleTimeout time.Duration

	// If wait is true and the pool is at the maxIdle limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	wait bool

	// @atomic
	closed uint32

	// Stack of idle Poolable with most recently used at the front.
	idle chan Poolable

	// Server address, e.g. "127.0.0.1:8080"
	// @const
	server string

	dpool *Pool

	// If marked as unavailable, then the checking goroutine will check it availability periodically.
	// A server is "available" if we can connnect to it, and respond to Ping() request of client.
	// Since no atomic boolean provided in Golang, we use uint32 instead.
	// @atomic
	available uint32

	// The failure count in succession. If the fails reached the threshold of "unavailable",
	// then this server should be marked as "unavailable", and we will not get connection
	// from it until recovered.
	// The idea of "fails" & "maxFails" is borrowed from Nginx.
	// @atomic
	fails uint32

	// The idea of "fails" & "maxFails" is borrowed from Nginx.
	// @const
	maxFails uint32

	stats PoolStats
}

// NewPoolShard creates a new pool shard.
func NewPoolShard(server string, parent *Pool, maxIdle, maxActive, maxFails int) *PoolShard {
	return &PoolShard{
		server:    server,
		dpool:     parent,
		maxIdle:   maxIdle,
		maxActive: int32(maxActive),
		idle:      make(chan Poolable, maxIdle),
		wait:      false, // TODO timed wait
		available: 1,
		closed:    0,
		maxFails:  uint32(maxFails),
	}
}

// If state changed, return true
func (p *PoolShard) markAvailable(b bool) (changed bool) {
	if b {
		return atomic.CompareAndSwapUint32(&p.available, 0, 1)
	}
	p.empty()
	return atomic.CompareAndSwapUint32(&p.available, 1, 0)
}

func (p *PoolShard) isAvailable() bool {
	if atomic.LoadUint32(&p.available) != 0 {
		return true
	}
	return false
}

// markFailed mark this shard as "suspect" or not
func (p *PoolShard) markFailed(failed bool) {
	if failed {
		s := atomic.AddUint32(&p.fails, 1)
		if s == p.maxFails {
			select {
			case p.dpool.suspectShards <- p:
			default: /*do nothing*/
			}
		}
	} else {
		atomic.StoreUint32(&p.fails, 0)
	}
}

func (p *PoolShard) suspectable() bool {
	if atomic.LoadUint32(&p.fails) >= p.maxFails {
		return true
	}
	return false
}

// Close releases the resources used by the pool shard.
func (p *PoolShard) Close() error {
	if !atomic.CompareAndSwapUint32(&p.closed, 0, 1) {
		return errors.New("dpool: close pool shard already closed")
	}
	p.empty()
	return nil
}

// get prunes stale connections and returns a connection from the idle channel or
// creates a new connection.
// The application must return the borrowed connection.
func (p *PoolShard) get() (c Poolable, err error) {
	// Check for pool closed before creating a new connection.
	if atomic.LoadUint32(&p.closed) == 1 {
		return nil, errors.New("dpool: get on closed pool shard")
	}

	atomic.AddUint64(&p.stats.NumGet, 1)

	select {
	case c = <-p.idle:
		c.setBorrowed(true)
		return c, nil
	default:
		if p.maxActive != 0 && atomic.LoadInt32(&p.active) >= p.maxActive {
			return nil, ErrPoolExhausted
		}

		// Dial new connection if under limit.
		atomic.AddInt32(&p.active, 1)

		atomic.AddUint64(&p.stats.NumDial, 1)
		if c, err = p.dpool.connFactory.Create(p.server); err != nil {
			p.markFailed(true)
			atomic.AddUint64(&p.stats.NumDialError, 1)
			atomic.AddInt32(&p.active, -1)
			return nil, err
		}

		// Setup pooled connection
		p.markFailed(false)
		c.setDataSource(p)
		c.setBorrowed(true)
		return c, nil
	}
}

func (p *PoolShard) put(c Poolable, broken bool) error {
	// XXX: Check if this object is borrowed and mark returning MUST be atomic
	c.lock()
	if !c.isBorrowed() {
		c.unlock()
		return ErrReturnInvalid
	}
	c.setBorrowed(false)
	c.unlock()

	p.markFailed(broken)
	atomic.AddUint64(&p.stats.NumPut, 1)
	if broken {
		atomic.AddUint64(&p.stats.NumBroken, 1)
	}

	if broken || atomic.LoadUint32(&p.closed) == 1 {
		atomic.AddUint64(&p.stats.NumClose, 1)
		atomic.AddInt32(&p.active, -1)
		return p.dpool.connFactory.Close(c)
	}

	select {
	case p.idle <- c:
	default:
		atomic.AddUint64(&p.stats.NumEvict, 1)
		atomic.AddUint64(&p.stats.NumClose, 1)
		atomic.AddInt32(&p.active, -1)
		return p.dpool.connFactory.Close(c)
	}

	return nil
}

// Empty removes and calls Close() on all the connections currently in the pool.
// Assuming there are no other connections waiting to be Put back this method
// effectively closes and cleans up the pool.
func (p *PoolShard) empty() {
	for {
		select {
		case c := <-p.idle:
			p.dpool.connFactory.Close(c)
			atomic.AddInt32(&p.active, -1)
			atomic.AddUint64(&p.stats.NumClose, 1)
		default:
			return
		}
	}
}

// Get statistics for this pool shard
func (p *PoolShard) getStats() (stats PoolStats) {
	stats.Shard = p.server
	stats.NumActive = int(atomic.LoadInt32(&p.active))
	if atomic.LoadUint32(&p.available) == 1 {
		stats.Available = true
	} else {
		stats.Available = false
	}

	stats.NumGet = atomic.SwapUint64(&p.stats.NumGet, 0)
	stats.NumPut = atomic.SwapUint64(&p.stats.NumPut, 0)
	stats.NumBroken = atomic.SwapUint64(&p.stats.NumBroken, 0)
	stats.NumClose = atomic.SwapUint64(&p.stats.NumClose, 0)
	stats.NumDial = atomic.SwapUint64(&p.stats.NumDial, 0)
	stats.NumDialError = atomic.SwapUint64(&p.stats.NumDialError, 0)
	stats.NumEvict = atomic.SwapUint64(&p.stats.NumEvict, 0)
	return
}
