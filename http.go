package zddgo

import(
	"net/http"
	"fmt"
	"time"
	"strconv"
	"github.com/feekk/zddgo/ztime"
)


func NewHttpSvr(c *HttpSvrConfig, handler http.Handler) (svr *http.Server){
	svr = &http.Server{
		Addr:           fmt.Sprintf(":%s", strconv.Itoa(c.Port)),
		Handler:        handler,
		ReadTimeout:    time.Duration(c.ReadTimeout),
		WriteTimeout:   time.Duration(c.WriteTimeout),
		MaxHeaderBytes: c.MaxHeaderBytes, //1 << 20 
	}
	return
}

type HttpSvrConfig struct{
	Port int
	ReadTimeout ztime.Duration
	WriteTimeout ztime.Duration
	MaxHeaderBytes int
}