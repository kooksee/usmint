package mint

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/types/consts"
)

func NewMetadata(dbs ... *kdb.KDB) *Miner {
	db1 := db
	if len(dbs) > 0 {
		db1 = dbs[0]
	}

	name := consts.Meta(consts.MetadataPrefix)
	return &Miner{name: name, db: db1.KHash([]byte(name))}
}

type Metadata struct {
	db   *kdb.KHash
	name string
}
