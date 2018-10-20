package validator

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn/consts"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	db kdb.IKHash
)

func Init(logger log.Logger) {
	db = kdb.DefaultConfig().GetDb().KHash([]byte(consts.ValidatorPrefix))
}
