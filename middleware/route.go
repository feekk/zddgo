package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
	"github.com/feekk/zddgo/errors"
	"github.com/feekk/zddgo/log"
	"github.com/feekk/zddgo/trace"
)

func Route() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		t := trace.InheritHttpTrace(ctx.Request)
		_, _, remoteAddr, _, _, _ := t.Get()
		ctx.Set(trace.TraceContextKey, t)

		//GetRawData return stream data.
		rawData, _ := ctx.GetRawData()
		//put back to Request.Body
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))

		//prepare route log parameter
		preLog := make(map[string]interface{})
		preLog["remote_ip"] = remoteAddr
		preLog["request_method"] = ctx.Request.Method
		preLog["url_path"] = ctx.Request.URL.Path
		if raw := ctx.Request.URL.RawQuery; raw != "" {
			preLog["url_path"] = preLog["url_path"].(string) + "?" + raw
		}
		preLog["raw_data"] = string(rawData)
		preLog["header"] = ctx.Request.Header
		log.Info(ctx, log.TAG_COM_REQUEST_IN, preLog)

		//for accept responese body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		ctx.Next()

		respLog := make(map[string]interface{})
		respLog["remote_ip"] = preLog["remote_ip"]
		respLog["request_method"] = preLog["request_method"]
		respLog["url_path"] = preLog["url_path"]
		respLog["raw_data"] = preLog["raw_data"]
		respLog["http_status"] = ctx.Writer.Status()
		respLog["http_message"] = ctx.Errors.ByType(gin.ErrorTypePrivate).String()
		respLog["resp_message"] = blw.body.String()
		respLog["cost_time_us"] = time.Now().Sub(start).Microseconds()
		log.Info(ctx, log.TAG_COM_REQUEST_OUT, respLog)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (i int, err error) {
	w.body.Write(b)
	i, err = w.ResponseWriter.Write(b)
	err = errors.With(err, "bodyLogWriter.Write.fail")
	return
}
