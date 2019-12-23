package validator

import (
	"mime/multipart"
	"net/http"
	"reflect"
	"github.com/gin-gonic/gin"
	"github.com/feekk/zddgo/errors"
)

const defaultMemory = 32 * 1024 * 1024

type formBinding struct{}

func (formBinding) Name() string {
	return "form"
}
func (formBinding) Tag() string {
	return "form"
}

func (b formBinding) Bind(ctx *gin.Context, obj interface{}) error {
	if err := ctx.Request.ParseForm(); err != nil {
		return err
	}
	if err := ctx.Request.ParseMultipartForm(defaultMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}
	rq, err := mapFormByTag(obj, ctx.Request.Form, b.Tag())
	_, ok := err.(VaildatorErrors)
	if err != nil && !ok {
		return errors.With(err)
	}
	err = vaildator(rq, obj, b.Tag(), err)
	return errors.With(err)
}

type formPostBinding struct{}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}
func (formPostBinding) Tag() string {
	return "form"
}

func (b formPostBinding) Bind(ctx *gin.Context, obj interface{}) error {
	if err := ctx.Request.ParseForm(); err != nil {
		return err
	}
	rq, err := mapFormByTag(obj, ctx.Request.PostForm, b.Tag())
	_, ok := err.(VaildatorErrors)
	if err != nil && !ok {
		return errors.With(err)
	}
	err = vaildator(rq, obj, b.Tag(), err)
	return errors.With(err)
}

type formMultipartBinding struct{}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}
func (formMultipartBinding) Tag() string {
	return "form"
}

func (b formMultipartBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}
	rq, err := mapSetterByTag(obj, (*multipartRequest)(req), b.Tag())
	_, ok := err.(VaildatorErrors)
	if err != nil && !ok {
		return errors.With(err)
	}
	err = vaildator(rq, obj, b.Tag(), err)
	return errors.With(err)
}

type multipartRequest http.Request

var _ setter = (*multipartRequest)(nil)

var (
	multipartFileHeaderStructType = reflect.TypeOf(multipart.FileHeader{})
)

func (r *multipartRequest) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (interface{}, error) {
	if value.Type() == multipartFileHeaderStructType {
		_, file, err := (*http.Request)(r).FormFile(key)
		if err != nil {
			return nil, err
		}
		if file != nil {
			value.Set(reflect.ValueOf(*file))
			return value.Interface(), nil
		}
	}
	return trySet(value, field, r.MultipartForm.Value, key, opt)
}