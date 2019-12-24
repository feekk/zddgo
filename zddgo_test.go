package zddgo

import(
	"flag"
	"runtime"
	"testing"
	"github.com/gin-gonic/gin"
)
func TestMain(t *testing.T){
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	z := New()

	if err := z.InitConfig(); err != nil{
		panic(err)
	}
	if err := z.InitOrm(); err != nil{
		panic(err)
	}
	if err := z.InitCache(); err != nil{
		panic(err)
	}
	engine := gin.New()
	//engine.Use(middleware.Recovery(), middleware.Route())
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "ping")
	})
	if err := z.HttpStart(engine); err != nil{
		panic(err)
	}
}