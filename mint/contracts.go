package mint

import (
	"github.com/yuin/gopher-lua"
	"github.com/kooksee/kdb"
	"github.com/layeh/gopher-luar"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kooksee/usmint/cmn"
	"github.com/patrickmn/go-cache"
	"github.com/kooksee/usmint/mint/luas"
)

type ContractManager struct {
	c  *cache.Cache
	db kdb.IKHash
}

func (c *ContractManager) getContract(addr []byte) *lua.LState {
	if ret, ok := c.c.Get(string(addr)); ok {
		return ret.(*lua.LState)
	}

	// 得到合约地址
	d, err := c.db.Get(addr)
	if err != nil {
		panic(err.Error())
	}

	l := lua.NewState()

	// 加载lua lib
	c.loadLib(l)

	cmn.MustNotErr("lua lib exec error", l.DoString(string(d)))

	c.c.SetDefault(string(addr), l)
	return l
}

func (c *ContractManager) loadLib(l *lua.LState) {
	//	加载系统类库
	l.SetGlobal("field", luar.New(l, c))

	// 加载数据库
	l.SetGlobal("field", luar.New(l, c))

	// 加载hash函数

}

func (c *ContractManager) Deploy() error {

	l := lua.NewState()
	defer l.Close()

	if err := l.DoString(""); err != nil {
		return err
	}

	err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("init"),
		NRet:    0,
		Protect: true}, lua.LString(""))
	if err != nil {
		return err
	}

	return nil
}

func (c *ContractManager) CallWithRet(cAddr []byte, method string, args []byte) ([]byte, error) {
	l := c.getContract(cAddr)
	val, err := luas.DecodeRaw(l, args)
	if err != nil {
		return nil, err
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal(method),
		NRet:    1,
		Protect: true,
	}, val); err != nil {
		return nil, err
	}

	ret := l.Get(-1)
	l.Pop(1)
	return luas.LValueDumps(ret)
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

func (c *Contract) CreateContractAddress(tx []byte) []byte {
	return crypto.Keccak256(tx, c.Code)
}
