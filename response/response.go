package response

import(
	"github.com/feekk/zddgo/validator"
	"github.com/gin-gonic/gin"
)

func JSONFail(ctx *gin.Context, data interface{}){
	r := _responseTempMap.get(FAIL)
	jsonTemp(ctx, r, data)
	return
}

func JSONSuccess(ctx *gin.Context, data interface{}){
	r := _responseTempMap.get(OK)
	jsonTemp(ctx, r, data)
	return
}

func JSONParamErr(ctx *gin.Context, err error){
	if err, ok := err.(validator.VaildatorErrors); !ok{
		JSONByCode(ctx, PARAMERR, map[string]interface{}{"data":err})
	}else{
		JSONByCode(ctx, PARAMERR, map[string]interface{}{"data":err.ToResponse()})
	}
	return
}

func JSONByCode(ctx *gin.Context, code int, data interface{}){
	r := _responseTempMap.get(code)
	jsonTemp(ctx, r, data)
	return
}

func jsonTemp(ctx *gin.Context, r responseTemp, data interface{}){
	ctx.JSON(r.status, h(r, data))
	return
}

func JSONCustomize(ctx *gin.Context, status, code int, message string, data interface{}){
	var r responseTemp = responseTemp{
		status : status,
		code : code,
		message : message,
	}
	ctx.JSON(status, h(r, data))
	return
}

func h(r responseTemp, data interface{}) (gh gin.H){
	gh = gin.H{
		"code": r.code,
		"message": r.message,
		"data":data,
	}
	return
}