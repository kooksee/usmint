package mint

import (
	"github.com/yuin/gopher-lua"
)

type Contract struct {
	Address []byte `json:"addr,omitempty"`
	Method  string `json:"method,omitempty"`
	Data    []byte `json:"data,omitempty"`
}

func (c *Contract) Deploy() {

	l := lua.NewState()
	defer l.Close()

	if err := l.DoString(""); err != nil {
		panic(err.Error())
	}

	err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("init"),
		NRet:    0,
		Protect: true}, lua.LString(""))
	if err != nil {
		panic(err.Error())
	}
}
