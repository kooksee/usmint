package validator

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
	db = kdb.DefaultConfig().GetDb().KHash([]byte("validator"))
	logger = cmn.Log().With("pkg", "validator")
}
