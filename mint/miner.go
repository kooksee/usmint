package mint

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/types/consts"
)

// 矿工设置

func NewMiner(dbs ... *kdb.KDB) *Miner {
	db1 := db
	if len(dbs) > 0 {
		db1 = dbs[0]
	}

	name := consts.Meta(consts.MinerPrefix)
	return &Miner{name: name, db: db1.KHash([]byte(name))}
}

type Miner struct {
	db   *kdb.KHash
	name string
}

// 设置矿工地址
// v 验证节点
// mn 矿工地址
func (m Miner) Set(v []byte, mn []byte) error {
	return m.db.Set(v, mn)
}

// 删除矿工地址
// v 验证节点
func (m Miner) Delete(v []byte) error {
	return m.db.Del(v)
}

// 记录矿工的交易地址
// 记录矿工的交易数量
// 计算矿工受益