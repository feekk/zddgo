package zddgo

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/feekk/zddgo/errors"
)

var (
	path string
	Conf *Config
)

func NewConfig(path string) (c *Config, err error) {
	c = &Config{}
	_, err = toml.DecodeFile(path, &c)
	err = errors.With(err)
	return
}

type Config struct {
	App      AppConf
	Env      EnvConf
	Http     HttpSvrConfig
	Logger   LoggerConf
	Trace    TraceConf
	Database DatabaseConf
	Redis    RedisConf
}

type AppConf struct {
	Name string
}

type EnvConf struct {
	GinMode   string
	DbLogMode bool
}

type LoggerConf struct{}

type TraceConf struct{}

type DatabaseConf struct {
	Use bool
	Default    OrmPoolConfig
	Connection []OrmPoolConfig
}

type RedisConf struct {
	Use bool
	Default    RedisPoolConfig
	Connection []RedisPoolConfig
}

func ConfigInit() (err error) {
	if path == "" {
		err = errors.New("empty config path")
	} else {
		Conf, err = NewConfig(path)
	}
	return
}

func init() {
	flag.StringVar(&path, "config", "", "default config path")
}
