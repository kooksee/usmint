package luas

import (
	"github.com/yuin/gopher-lua"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"math/big"
)

func NewDb(contractAddress []byte) *ContractDb {
	return &ContractDb{db: db.KHash(contractAddress)}
}

type ContractDb struct {
	address []byte
	l       *lua.LState
	db      kdb.IKHash
}

func (c *ContractDb) Db(name string) *CDb {
	return &CDb{db: db.KHash(append(c.address, []byte(name)...)), l: c.l}
}

type CDb struct {
	db kdb.IKHash
	l  *lua.LState
}

func (c *CDb) SetStr(key, value string) {
	cmn.ErrPipe("CDb SetStr Error", c.db.WithBatch(func(kh kdb.IKHBatch) error {
		return kh.Set([]byte(key), []byte(value))
	}))
}

func (c *CDb) Str(key string) string {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipe("CDb Str Error", err); err != nil {
		return ""
	}
	return string(val)
}

func (c *CDb) SetInt(key string, value int) {
	cmn.ErrPipe("CDb SetInt Error", c.db.WithBatch(func(kh kdb.IKHBatch) error {
		return kh.Set([]byte(key), big.NewInt(int64(value)).Bytes())
	}))
}

func (c *CDb) Int(key string) int {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipe("CDb Int Error", err); err != nil {
		return 0
	}
	return int(big.NewInt(0).SetBytes(val).Int64())
}

func (c *CDb) SetTable(key string, m map[string]interface{}) {
	c.db.WithBatch(func(kh kdb.IKHBatch) error {
		val, err := json.Marshal(m)
		return cmn.ErrPipe("CDb SetTable Error", err, kh.Set([]byte(key), val))
	})
}

func (c *CDb) Table(key string) lua.LValue {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipe("CDb Table Get Error", err); err != nil {
		return c.l.NewTable()
	}
	v, err := decodeRaw(c.l, val)
	if err := cmn.ErrPipe("CDb Table decodeRaw Error", err); err != nil {
		return c.l.NewTable()
	}
	return v
}
