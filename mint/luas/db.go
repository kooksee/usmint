package luas

import (
	"github.com/yuin/gopher-lua"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
)

func NewDb(contractAddress []byte) *ContractDb {
	return &ContractDb{db: db.KHash(contractAddress)}
}

type ContractDb struct {
	address []byte
	l       *lua.LState
	db      *kdb.KHash
}

func (c *ContractDb) NewT(name string) *kdb.KHash {
	return db.KHash(append(c.address, []byte(name)...))
}

func (c *ContractDb) SetStr(key, value string) {
	cmn.ErrPipeLog("ContractDb SetStr Error", c.db.WithTx(func(kh *kdb.KHBatch) error {
		return kh.Set([]byte(key), []byte(value))
	}))
}

func (c *ContractDb) Str(key string) string {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipeLog("ContractDb Str Error", err); err != nil {
		return ""
	}
	return string(val)
}

func (c *ContractDb) SetInt(key string, value int) {
	c.db.WithTx(func(kh *kdb.KHBatch) error {
		return kh.Set([]byte(key), cmn.Int64ToByte(int64(value)))
	})
}

func (c *ContractDb) Int(key string) int {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipeLog("ContractDb Int Error", err); err != nil {
		return 0
	}
	return int(cmn.ByteToInt64(val))
}

func (c *ContractDb) SetFloat(key string, value float64) {
	cmn.ErrPipeLog("ContractDb SetFloat Error", c.db.WithTx(func(kh *kdb.KHBatch) error {
		return kh.Set([]byte(key), cmn.Float64ToByte(value))
	}))
}

func (c *ContractDb) Float(key string) float64 {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipeLog("ContractDb Float Error", err); err != nil {
		return 0
	}
	return cmn.ByteToFloat64(val)
}

func (c *ContractDb) SetTable(key string, m map[string]interface{}) {
	cmn.ErrPipeLog("ContractDb SetTable Error", c.db.WithTx(func(kh *kdb.KHBatch) error {
		val, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return kh.Set([]byte(key), val)
	}))
}

func (c *ContractDb) Table(key string) lua.LValue {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipeLog("ContractDb Table Get Error", err); err != nil {
		return c.l.NewTable()
	}
	v, err := decodeRaw(c.l, val)
	if err := cmn.ErrPipeLog("ContractDb Table decodeRaw Error", err); err != nil {
		return c.l.NewTable()
	}
	return v
}

// 自定义全局类型
// 自定义数据存储
// 自定义kv存储
// 合约编号,所有的合约按照递增序列排列

/*
1. 部署合约
2. 获得合约地址,获得合约编号,合约编号前缀
3. 合约数据存储，contracts:1234:name12345:ssss:yuuee value

合约元数据存储,合约的code，编号，发送者，公钥，以及nonce
根据合约地址获取合约的code，获得合约的编号
meta:address {code,编号,其他的元数据}
静态数据放入数据库存储，合约数据用merker tree去存储，然后得到app hash
不管中间逻辑多少，只要是结果是一样的，那就没问题
4. 静态数据存储，code，编号，区块信息，tx
 */
