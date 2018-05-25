package code

import "fmt"

func errs(code uint32, msg string) *E {
	return &E{
		Code:code,
		Msg:msg,
	}
}

type E struct {
	Code uint32
	Msg  string
}

func (e *E) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

