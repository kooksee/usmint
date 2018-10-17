package validator

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn/consts"
)

var (
	db kdb.IKHash
)

func Init() {
	db = kdb.DefaultConfig().GetDb().KHash([]byte(consts.ValidatorPrefix))
}
