package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/feekk/zddgo/errors"
)

type uriBinding struct{}

func (uriBinding) Name() string {
	return "uri"
}

func (uriBinding) Tag() string{
	return "uri"
}

func (b uriBinding) Bind(ctx *gin.Context, obj interface{}) error {
	m := make(map[string][]string)
	for _, v := range ctx.Params {
		m[v.Key] = []string{v.Value}
	}
	rq, err := mapFormByTag(obj, m, b.Tag())
	_, ok := err.(VaildatorErrors)
	if err != nil && !ok {
		return errors.With(err)
	}
	err = vaildator(rq, obj, b.Tag(), err)
	return errors.With(err)
}
