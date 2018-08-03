package code

import "fmt"

func errs(code uint32, msg string) *e {
	return &e{
		Code: code,
		Msg:  msg,
	}
}

type e struct {
	Code uint32
	Msg  string
}

func (e *e) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Msg)
}
