package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"github.com/feekk/zddgo/log"
	"github.com/feekk/zddgo/trace"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				request, _ := httputil.DumpRequest(ctx.Request, true)
				stack := trace.Stack(3)
				stackByte, _ := json.Marshal(fmt.Sprintf("%s %s", string(request), string(stack)))

				if brokenPipe {
					log.Warn(ctx, log.TAG_Status_Internal_Server_BrokenPipe, map[string]interface{}{
						"err":   err,
						"stack": string(stackByte),
					})
				} else {
					log.Error(ctx, log.TAG_Status_Internal_Server_Error, map[string]interface{}{
						"err":   err,
						"stack": string(stackByte),
					})
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					ctx.Error(err.(error)) // nolint: errcheck
					ctx.Abort()
				} else {
					//ctx.Abort()
					//response.ResponseWithHttpCode(ctx, http.StatusInternalServerError, response.StatusSystemError,nil)
					ctx.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()

		ctx.Next()
	}
}
