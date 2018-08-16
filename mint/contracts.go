package mint

import (
	"github.com/yuin/gopher-lua"
	"github.com/kooksee/kdb"
	"github.com/layeh/gopher-luar"
	"github.com/kooksee/usmint/cmn"
	"github.com/patrickmn/go-cache"
	"github.com/kooksee/usmint/mint/luas"
	"time"
	"github.com/kooksee/usmint/kts/consts"
)

func newContractManager() *ContractManager {
	return &ContractManager{
		c:  cache.New(time.Hour, time.Hour),
		db: db.KHash([]byte(consts.SmartContractPrefix)),
	}
}

type ContractManager struct {
	c  *cache.Cache
	db kdb.IKHash
}

func (c *ContractManager) getContract(addr []byte) *lua.LState {
	if ret, ok := c.c.Get(string(addr)); ok {
		return ret.(*lua.LState)
	}

	// 得到合约地址
	d, err := c.db.Get(append(addr, "code"...))
	cmn.MustNotErr("ContractManager getContract 1", err)

	l := lua.NewState()
	cmn.MustNotErr("ContractManager getContract 2", l.DoString(string(d)))

	di, err := c.db.Get(append(addr, "init"...))
	cmn.MustNotErr("ContractManager getContract 3", err)

	lv, err := luas.DecodeRaw(l, di)
	cmn.MustNotErr("ContractManager getContract 4", err)

	// 加载init
	l.SetGlobal("init", lv)

	// 加载lua lib
	c.loadLib(l)

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

// Deploy 部署合约
func (c *ContractManager) Deploy(address []byte, data string) error {

	l := lua.NewState()
	defer l.Close()

	if err := l.DoString(data); err != nil {
		return cmn.ErrPipe("ContractManager Deploy 1", err)
	}

	t, ok := l.GetGlobal("init").(*lua.LTable);
	if !ok {
		return cmn.Err("ContractManager DeployCheck 2: the init variable must be table type")
	}

	// 保存合约
	if err := c.db.Set(append(address, "code"...), []byte(data)); err != nil {
		return cmn.ErrPipe("ContractManager DeployCheck 3", err)
	}

	// 保存合约的init
	dt, err := luas.LValueDumps(t)
	if err != nil {
		return cmn.ErrPipe("ContractManager DeployCheck 4", err)
	}

	return cmn.ErrPipe("ContractManager DeployCheck 5", c.db.Set(append(address, "init"...), dt))
}

// DeployCheck 合约部署检查
func (c *ContractManager) DeployCheck(addr []byte, data string) error {
	// 检查合约地址是否存在

	ok, err := c.db.Exist(addr)
	if err != nil {
		return cmn.ErrPipe("ContractManager DeployCheck", err)
	}

	if ok {
		return cmn.Err("ContractManager DeployCheck: the contract %s had exist", addr)
	}

	l := lua.NewState()
	defer l.Close()

	if err := l.DoString(data); err != nil {
		return cmn.ErrPipe("ContractManager DeployCheck", err)
	}

	if _, ok := l.GetGlobal("init").(*lua.LTable); !ok {
		return cmn.Err("ContractManager DeployCheck: the init variable must be table type")
	}

	return nil
}

func (c *ContractManager) CallWithRetCheck(cAddr []byte, method string, args []byte) error {
	if method == "init" {
		return cmn.Err("ContractManager CallWithRetCheck: the method %s is keyword", method)
	}

	l := c.getContract(cAddr)
	_, err := luas.DecodeRaw(l, args)
	return cmn.ErrPipe("ContractManager CallWithRetCheck", err)
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
