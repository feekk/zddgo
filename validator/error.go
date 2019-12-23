package validator

import (
	"bytes"
	"fmt"
	"strings"
)

type VaildatorErrors map[string]error

func (z VaildatorErrors) Error() string {
	buff := bytes.NewBufferString("")
	for field, err := range z {
		buff.WriteString(fmt.Sprintf("field:%s error:%s", field, err))
		buff.WriteString(" ")
	}
	return strings.TrimSpace(buff.String())
}
func (z VaildatorErrors) ToResponse() map[string]string {
	resp := make(map[string]string)
	for field, err := range z {
		resp[field] = err.Error()
	}
	return resp
}