package validator

import (
	"bytes"
	"io"
	"net/http"
	"encoding/json"
	"github.com/feekk/zddgo/errors"
)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Tag() string {
	return "json"
}

func (b jsonBinding) Bind(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}
	if err := decodeJSON(req.Body, obj); err != nil {
		return errors.With(err)
	}
	rq, err := mapNil(obj)
	return errors.With(vaildator(rq, obj, b.Tag(), err))
}

func (jsonBinding) BindBody(body []byte, obj interface{}) error {
	if err := decodeJSON(bytes.NewReader(body), obj); err != nil {
		return errors.With(err)
	}
	_, err := mapNil(obj)
	return errors.With(err)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(obj); err != nil {
		return errors.With(err)
	}
	return nil
}