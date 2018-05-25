package app

import (
	"errors"
	"github.com/yuin/gopher-lua"
	"github.com/hashicorp/golang-lru/simplelru"
	"github.com/tendermint/tmlibs/common"
	"github.com/layeh/gopher-luar"
)

type ContractManager struct {
	cache simplelru.LRUCache
	size  int
}

func NewContractManager(size int) *ContractManager {
	l, err := simplelru.NewLRU(size, nil)
	if err != nil {
		panic(err.Error())
	}
	return &ContractManager{
		cache: l,
		size:  size,
	}
}

func (c *ContractManager) AddContract(addr string) error {
	code := state.db.Get([]byte(addr))
	if code == nil {
		return errors.New("合约代码不存在")
	}

	l := lua.NewState()
	l.SetGlobal("k", luar.New(l, state))
	if err := l.DoString(string(code)); err != nil {
		return err
	}

	if c.cache.Len() > c.size {
		k, _, ok := c.cache.GetOldest()
		if ok {
			c.DelContract(common.Fmt("%s", k))
		}
	}

	c.cache.Add(addr, l)
	return nil
}

func (c *ContractManager) GetContract(addr string) *lua.LState {

	if v, ok := c.cache.Get(addr); ok {
		return v.(*lua.LState)
	}

	return nil
}

func (c *ContractManager) DelContract(addr string) bool {
	ctt := c.GetContract(addr)
	if ctt != nil {
		ctt.Close()
	}
	return c.cache.Remove(addr)
}

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
		Fn:      l.GetGlobal(c.Method),
		NRet:    0,
		Protect: true}, lua.LString(""))
	if err != nil {
		panic(err.Error())
	}

}
