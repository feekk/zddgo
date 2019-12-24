package zddgo

import(
	"io"
	"net/http"
	"fmt"
	"time"
	"strconv"
	"context"
	"github.com/feekk/zddgo/ztime"
	"github.com/feekk/zddgo/trace"
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

//
//
//
func NewHttpRequest(ctx context.Context) (c *HttpRequest){
	c = &HttpRequest{
		ctx : ctx,
		client : &http.Client{},
	}
	return
}
type HttpRequest struct{
	ctx context.Context
	client *http.Client
}
func(r *HttpRequest) Get(url string)(resp *http.Response, err error){
	req, err := http.NewRequest(http.MethodGet, url, nil)
	return r.do(req)
}
func(r *HttpRequest) Post(url, contentType string, body io.Reader)(resp *http.Response, err error){
	req, err := http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	return r.do(req)
}
func(r *HttpRequest) do(req *http.Request) (resp *http.Response, err error){
	t := trace.InheritContextTrace(r.ctx)
	if t != nil {
		t.IncrRpc()
		traceId, spanId, _, _, _, _ := t.Get()
		req.Header.Set(trace.GetTraceHeadKey(), traceId)
		req.Header.Set(trace.GetSpandHeadKey(), spanId)
	}
	return r.client.Do(req)
}