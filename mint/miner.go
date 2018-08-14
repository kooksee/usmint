package mint

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/kts/consts"
	"github.com/kooksee/usmint/cmn"
)

// 矿工设置
func NewMiner() *Miner {
	name := consts.Meta(consts.MinerPrefix)
	return &Miner{name: name, db: db.KHash([]byte(name))}
}

type Miner struct {
	db   kdb.IKHash
	name string
}

// 设置矿工地址
// v 验证节点
// mn 矿工地址
func (m Miner) Set(v []byte, mn []byte) error {
	return cmn.ErrPipe("Miner.Set", m.db.Set(v, mn))
}

// 删除矿工地址
// v 验证节点
func (m Miner) Delete(v []byte) error {
	return cmn.ErrPipe("Miner.Delete", m.db.Del(v))
}
