package zddgo

import(
	"sync"
	"github.com/jinzhu/gorm"
)

var (
	mysqlPool = &MysqlConns{conn:&sync.Map{}}
)

func MysqlInit(c *DatabaseConf) (err error){
	if !c.Use {
		return
	}
	
	if mysqlPool.def, err = NewOrmPool(&c.Default); err != nil {
		return
	}
	
	if len(c.Connection) <= 0 {
		return
	}
	var conn *gorm.DB
	for _, conConf := range c.Connection{
		conn, err = NewOrmPool(&conConf)
		if err != nil{
			return 
		}
		mysqlPool.conn.Store(conConf.Name, *conn)
	} 
	return
}

type MysqlConns struct {
	def  *gorm.DB
	conn  *sync.Map
}

//
//
func MysqlDef() *gorm.DB{
	return mysqlPool.def
}
//
//
func MysqlConn(key string) (*gorm.DB, bool){
	if pool, ok := mysqlPool.conn.Load(key); ok {
		 p := pool.(gorm.DB)
		return &p, true
	}
	return nil, false
}