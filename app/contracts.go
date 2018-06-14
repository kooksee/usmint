package app

import (
	"errors"
	"github.com/yuin/gopher-lua"
	"github.com/layeh/gopher-luar"
	"github.com/patrickmn/go-cache"
	"time"
)

type ContractManager struct {
	cache *cache.Cache
}

func NewContractManager() *ContractManager {
	return &ContractManager{
		cache: cache.New(5*time.Minute, 10*time.Minute),
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

	return c.cache.Add(addr, l, cache.DefaultExpiration)
}

func (c *ContractManager) GetContract(addr string) *lua.LState {

	if v, ok := c.cache.Get(addr); ok {
		return v.(*lua.LState)
	}

	return nil
}

func (c *ContractManager) DelContract(addr string) {
	if ctt := c.GetContract(addr); ctt != nil {
		ctt.Close()
		c.cache.Delete(addr)
	}

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
