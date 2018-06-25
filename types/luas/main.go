package luas

import (
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"github.com/tendermint/go-crypto"
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

func (c *Contract) LoadLib() {
	//	加载系统类库
	c.l.SetGlobal("field", luar.New(c.l, c))
}
