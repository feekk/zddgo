package errors

import(
	"testing"
	"errors"
)

func TestNew(t *testing.T){
	err1 := errors.New("errors.New")
	err2 := New("github.com/pkg/errors")
	if errors.Is(err2, err1) || errors.Is(err2, err1){
		t.Errorf("github.com/pkg/errors\n")
	}
}

