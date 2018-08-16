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
	"github.com/kooksee/usmint/kts"
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
	l.SetGlobal("field", luar.New(l, c))
}

// Deploy 部署合约
func (c *ContractManager) DeployWithTx(tx *kts.Transaction) error {

	// 验证签名
	addr, err := tx.VerifySign()
	if err != nil {
		return err
	}

	// 得到合约的地址
	scAddr := kts.CreateContractAddress(addr.Bytes(), []byte(tx.Data))

	l := lua.NewState()
	defer l.Close()

	if err := l.DoString(tx.Data); err != nil {
		return cmn.ErrPipe("ContractManager Deploy 1", err)
	}

	t, ok := l.GetGlobal("init").(*lua.LTable);
	if !ok {
		return cmn.Err("ContractManager DeployCheck 2: the init variable must be table type")
	}

	// 保存合约
	if err := c.db.Set(append(scAddr.Bytes(), "code"...), []byte(tx.Data)); err != nil {
		return cmn.ErrPipe("ContractManager DeployCheck 3", err)
	}

	// 保存合约的init
	dt, err := luas.LValueDumps(t)
	if err != nil {
		return cmn.ErrPipe("ContractManager DeployCheck 4", err)
	}

	return cmn.ErrPipe("ContractManager DeployCheck 5", c.db.Set(append(scAddr.Bytes(), "init"...), dt))
}

// DeployCheckWithTx 合约部署检查
func (c *ContractManager) DeployCheckWithTx(tx *kts.Transaction) error {
	// 验证签名
	addr, err := tx.VerifySign()
	if err != nil {
		return err
	}

	// 得到合约的地址
	scAddr := kts.CreateContractAddress(addr.Bytes(), []byte(tx.Data))
	// 检查合约地址是否存在

	ok, err := c.db.Exist(scAddr.Bytes())
	if err != nil {
		return cmn.ErrPipe("ContractManager DeployCheck 1", err)
	}

	if ok {
		return cmn.Err("ContractManager DeployCheck: the contract %s had exist", addr)
	}

	l := lua.NewState()
	defer l.Close()

	if err := l.DoString(tx.Data); err != nil {
		return cmn.ErrPipe("ContractManager DeployCheck 2", err)
	}

	if _, ok := l.GetGlobal("init").(*lua.LTable); !ok {
		return cmn.Err("ContractManager DeployCheck: the init variable must be table type")
	}

	return nil
}

// CallWithRetCheck 合约调用检查
func (c *ContractManager) CallCheckWithTx(tx *kts.Transaction) error {
	addr, err := tx.VerifySign()
	if err != nil {
		return err
	}

	cnt, err := kts.DecodeContract(tx.Data)
	if err != nil {
		return err
	}

	if cnt.Method == "init" {
		return cmn.Err("ContractManager CallWithRetCheck: the method %s is keyword", cnt.Method)
	}

	_, err = luas.DecodeRaw(c.getContract(addr.Bytes()), cnt.Data)
	return cmn.ErrPipe("ContractManager CallWithRetCheck", err)
}

// CallWithOutRet 合约写入调用
func (c *ContractManager) CallWithOutRetWithTx(tx *kts.Transaction) error {
	// 验证签名
	addr, err := tx.VerifySign()
	if err != nil {
		return err
	}

	cnt, err := kts.DecodeContract(tx.Data)
	if err != nil {
		return err
	}

	l := c.getContract(cnt.Address)
	// 设置调用者
	l.SetGlobal("msg_sender", lua.LString(addr.Hex()))

	val, err := luas.DecodeRaw(c.getContract(cnt.Address), cnt.Data)
	return cmn.ErrPipe("ContractManager CallWithOutRet", err, cmn.ErrCurry(l.CallByParam, lua.P{
		Fn:      l.GetGlobal(cnt.Method),
		NRet:    0,
		Protect: true,
	}, val))
}

// CallWithRet 合约查询调用
func (c *ContractManager) CallWithRet(cAddr []byte, method string, args []byte) ([]byte, error) {
	l := c.getContract(cAddr)
	val, err := luas.DecodeRaw(l, args)
	if err != nil {
		return nil, cmn.ErrPipe("ContractManager CallWithRet 1", err)
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal(method),
		NRet:    1,
		Protect: true,
	}, val); err != nil {
		return nil, cmn.ErrPipe("ContractManager CallWithRet 2", err)
	}

	ret := l.Get(-1)
	l.Pop(1)
	return luas.LValueDumps(ret)
}
