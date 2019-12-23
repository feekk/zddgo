package zddgo

import(
	"time"
	"github.com/feekk/zddgo/ztime"
	"github.com/feekk/zddgo/errors"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

func NewOrmPool(c *OrmPoolConfig) (db *gorm.DB, err error){
	db, err = gorm.Open(c.Driver, c.DSN)
	if err != nil {
		err = errors.With(err)
		return 
	}
	db.DB().SetMaxIdleConns(c.MaxIdle)
	db.DB().SetMaxOpenConns(c.MaxConn)
	db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout))
	db.LogMode(c.LogMode)
	if err = db.DB().Ping(); err != nil{
		err = errors.With(err)
	}
	return 
}

type OrmPoolConfig struct{
	LogMode bool
	Driver string
	Name string
	DSN string         // data source name.
	MaxIdle int
	MaxConn int
	IdleTimeout ztime.Duration //Max life time 
}