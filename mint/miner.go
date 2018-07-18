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
func (m Miner) Set(v []byte, mn []byte) error {
	//	验证节点和矿工地址
	return m.db.Set(v, mn)
}

// 删除矿工地址
func (m Miner) Delete(v []byte) error {
	//	验证节点和矿工地址

	return m.db.Del(v)
}
