package zddgo

import(
	"testing"
	"github.com/gin-gonic/gin"
)

func TestZddgo(t *testing.T){
	z := Zddgo{}
	if err := z.InitConfig(); err != nil{
		panic(err)
	}
	if err := z.InitOrm(); err != nil{
		panic(err)
	}
	if err := z.InitCache(); err != nil{
		panic(err)
	}

	//gin.SetMode(Conf.EnvConf.GinMode)
	engine := gin.New()
	//engine.Use(middleware.Recovery(), middleware.Route())
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "ping")
	})
	if err := z.HttpStart(engine); err != nil{
		panic(err)
	}
}