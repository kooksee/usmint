package state

import (
	"github.com/tendermint/tendermint/libs/log"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
)

var (
	db     kdb.IKHash
	logger log.Logger
)

func Init() {
	cfg := kdb.DefaultConfig()
	db = cfg.GetDb().KHash([]byte("state"))
	logger = cmn.Log().With("pkg", "state")
}
