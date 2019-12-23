package zddgo

import (
	"fmt"
	"os"
	"strings"
	"time"
	"context"
)

//public tag
const(
	//http request in
	TAG_COM_REQUEST_IN = "_com_request_in"
	//http request out
	TAG_COM_REQUEST_OUT = "_com_request_out"
	//
	TAG_Status_Internal_Server_BrokenPipe = "_status_internal_server_brokenpipe"
	// 
	TAG_Status_Internal_Server_Error = "_status_internal_server_error"
)

var (
	symbol []string =[]string{
		"\t",
		"\n",
	} 
)

//[INFO][2018-12-12 00:00:00][main.go(32)] tag||trace_id=xxx||SpanId=XXX||RpcId=XXX||map[string][]interface{}
//var Formatter string = "[%s][%s][%s] %s||TraceId=%s||SpanId=%s||RpcId=%d%s\r\n"
var Formatter string = "[%s][%s] %s||TraceId=%s||SpanId=%s||RpcId=%d%s\r\n"

type logger struct{}

func(l *logger) print(ctx context.Context, level, tag string, parameter map[string]interface{}){
	var traceId, spanId, content string
	var rpcId int
	//trace
	if t := InheritContextTrace(ctx); t != nil{
		traceId, spanId, _, _, _, rpcId = t.Get()
	}
	// parameter
	if parameter != nil {
		var keyslice []string
		for idx, val := range parameter {
			keyslice = append(keyslice, fmt.Sprintf("||%s=%+v", idx, val))
		}
		content = strings.Replace(strings.Trim(fmt.Sprint(keyslice), "[]"), " ||", "||", -1)
	}
	for i:=0; i < len(symbol); i++ {
		content = strings.ReplaceAll(content, symbol[i], "")
	}
	//handle
	l.handler(
		Formatter, 
		strings.ToUpper(level), time.Now().Format("2006/01/02 15:04:05.000"), 
		//fmt.Sprintf("%s(%d)", file, line, 
		tag, traceId, spanId, rpcId, content,
	)
}
func(l *logger) handler(f string, v ...interface{}){
	fmt.Fprintf(os.Stdout, f, v...)
}


var(
	defaultLog = &logger{}
)


func Info(ctx context.Context, tag string, parameter map[string]interface{}){
	defaultLog.print(ctx, "INFO", tag, parameter)
}
func Warn(ctx context.Context, tag string, parameter map[string]interface{}){
	defaultLog.print(ctx, "WARNING", tag, parameter)
}
func Error(ctx context.Context, tag string, parameter map[string]interface{}){
	defaultLog.print(ctx, "ERROR", tag, parameter)
}
