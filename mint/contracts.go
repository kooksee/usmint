package mint

import (
	"github.com/yuin/gopher-lua"
	"github.com/kooksee/kdb"
	"github.com/layeh/gopher-luar"
	"github.com/kooksee/usmint/mint/luas"
	"github.com/ethereum/go-ethereum/crypto"
)

type ContractManager struct {
	c  map[string]*lua.LState
	db *kdb.KHash
}

func (c *ContractManager) Contract(addr []byte) *lua.LState {
	if ll, ok := c.c[string(addr)]; ok {
		return ll
	}

	d, err := c.db.Get(addr)
	if err != nil {
		panic(err.Error())
	}

	l := lua.NewState()
	if err := l.DoString(string(d)); err != nil {
		panic(err.Error())
	}

	c.c[string(addr)] = l
	return l
}

type Contract struct {
	Address         []byte `json:"addr,omitempty"`
	Method          string `json:"method,omitempty"`
	Data            []byte `json:"data,omitempty"`
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
	return crypto.Keccak256(tx, c.Code)
}

func (c *Contract) CallWithRet(method string, args []byte) ([]byte, error) {
	val, err := luas.DecodeRaw(c.l, args)
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
	return luas.LValueDumps(ret)
}

func (c *Contract) LoadLib() {
	//	加载系统类库
	c.l.SetGlobal("field", luar.New(c.l, c))

	// 加载数据库
	c.l.SetGlobal("field", luar.New(c.l, c))
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
