package errors

import (
	"errors"
	perror "github.com/pkg/errors"
)

func New(msg string) error {
	return perror.New(msg)
}

func With(err error, msg ...string) error{
	if err == nil {
		return nil
	}
	ferr := perror.New("")
	serr := perror.WithStack(errors.New(""))
	if errors.Is(err, ferr) || errors.Is(err, serr){
		return err
	}
	return perror.WithStack(err)
}

func Errorf(format string, args ...interface{}) error {
	return perror.Errorf(format, args...)
}