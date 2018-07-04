package luas

import (
	"github.com/yuin/gopher-lua"
	"github.com/tendermint/go-crypto"
	"github.com/layeh/gopher-luar"
)

func NewContract() *Contract {
	return &Contract{debug: false, l: lua.NewState()}
}

type Contract struct {
	l               *lua.LState
	debug           bool
	ContractAddress []byte
	Code            []byte
	Tx              []byte
}

func (c *Contract) Init() error {
	return nil
}

func (c *Contract) CreateContractAddress(tx []byte) []byte {
	return crypto.Ripemd160(append(tx, crypto.Ripemd160(c.Code)...))
}

func (c *Contract) Call(method string, args string) error {
	val, err := decodeRaw(c.l, []byte(args))
	if err != nil {
		return err
	}

	return c.l.CallByParam(lua.P{
		Fn:      c.l.GetGlobal(method),
		NRet:    0,
		Protect: true,
	}, val)
}

func (c *Contract) CallWithRet(method string, args string) ([]byte, error) {
	val, err := decodeRaw(c.l, []byte(args))
	if err != nil {
		return nil, err
	}

	if err := c.l.CallByParam(lua.P{
		Fn:      c.l.GetGlobal(method),
		NRet:    1,
		Protect: true,
	}, val); err != nil {
		return nil, err
	}

	ret := c.l.Get(-1)
	c.l.Pop(1)
	return LValueDumps(ret)
}

func (c *Contract) LoadLib() {
	//	加载系统类库
	c.l.SetGlobal("field", luar.New(c.l, c))

	// 加载数据库

	c.l.SetGlobal("field", luar.New(c.l, c))

}
