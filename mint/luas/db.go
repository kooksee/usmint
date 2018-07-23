package luas

import (
	"github.com/yuin/gopher-lua"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
)

// db:map
// db:set
// db:list

func NewDb(contractAddress []byte) *ContractDb {
	return &ContractDb{db: db.KHash(contractAddress)}
}

type ContractDb struct {
	address []byte
	l       *lua.LState
	db      *kdb.KHash
}

func (c *ContractDb) Db(name string) *CDb {
	return &CDb{db: db.KHash(append(c.address, []byte(name)...)), l: c.l}
}

type CDb struct {
	db *kdb.KHash
	l  *lua.LState
}

func (c *CDb) SetStr(key, value string) {
	cmn.ErrPipeLog("CDb SetStr Error", c.db.WithTx(func(kh *kdb.KHBatch) error {
		return kh.Set([]byte(key), []byte(value))
	}))
}

func (c *CDb) Str(key string) string {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipeLog("CDb Str Error", err); err != nil {
		return ""
	}
	return string(val)
}

func (c *CDb) SetInt(key string, value int) {
	cmn.ErrPipeLog("CDb SetInt Error", c.db.WithTx(func(kh *kdb.KHBatch) error {
		return kh.Set([]byte(key), cmn.Int64ToByte(int64(value)))
	}))
}

func (c *CDb) Int(key string) int {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipeLog("CDb Int Error", err); err != nil {
		return 0
	}
	return int(cmn.ByteToInt64(val))
}

// 该类型取消
func (c *CDb) SetFloat(key string, value float64) {
	cmn.ErrPipeLog("CDb SetFloat Error", c.db.WithTx(func(kh *kdb.KHBatch) error {
		return kh.Set([]byte(key), cmn.Float64ToByte(value))
	}))
}

// 该类型取消
func (c *CDb) Float(key string) float64 {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipeLog("CDb Float Error", err); err != nil {
		return 0
	}
	return cmn.ByteToFloat64(val)
}

func (c *CDb) SetTable(key string, m map[string]interface{}) {
	c.db.WithTx(func(kh *kdb.KHBatch) error {
		val, err := json.Marshal(m)
		return cmn.ErrPipeLog("CDb SetTable Error", err, kh.Set([]byte(key), val))
	})
}

func (c *CDb) Table(key string) lua.LValue {
	val, err := c.db.Get([]byte(key))
	if err := cmn.ErrPipeLog("CDb Table Get Error", err); err != nil {
		return c.l.NewTable()
	}
	v, err := decodeRaw(c.l, val)
	if err := cmn.ErrPipeLog("CDb Table decodeRaw Error", err); err != nil {
		return c.l.NewTable()
	}
	return v
}
