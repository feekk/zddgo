package utils


import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"github.com/gin-gonic/gin"
	"github.com/feekk/zddgo/log"
)

func StdPrint(format string, values ...interface{}) {
	fmt.Fprintf(os.Stdout, format, values...)
}

func GoWithRecover(ctx *gin.Context, f reflect.Value, args ...interface{}){
	defer func() {
		if err := recover(); err != nil {
			var buf [1024]byte
			n := runtime.Stack(buf[:], false)
			log.Error(ctx, log.TAG_Status_Internal_Server_Error, map[string]interface{}{
				"err":   err,
				"stack": string(buf[:n]),
			})
			
		}
	}()
	n := len(args)
	params := make([]reflect.Value, n)
	for i:=0; i < n; i++ {
		params[i] = reflect.ValueOf(args[i])
	}
	f.Call(params)
}