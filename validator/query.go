package validator

import (
	"github.com/feekk/zddgo/errors"
	"github.com/gin-gonic/gin"
)

type queryBinding struct{}

func (queryBinding) Name() string {
	return "query"
}
func (queryBinding) Tag() string{
	return "query"
}

func (b queryBinding) Bind(ctx *gin.Context, obj interface{}) error {
	values := ctx.Request.URL.Query()
	rq, err := mapFormByTag(obj, values, b.Tag())
	_, ok := err.(VaildatorErrors)
	if err != nil && !ok {
		return errors.With(err)
	}
	err = vaildator(rq, obj, b.Tag(), err)
	return errors.With(err)
}
