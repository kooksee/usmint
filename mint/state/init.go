package state

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn/consts"
	"github.com/kooksee/usmint/cmn/wire"
)

var (
	db kdb.IKHash
)

func Init() {
	cfg := kdb.DefaultConfig()
	db = cfg.GetDb().KHash([]byte(consts.StatePrefix))
}

func init() {
	wire.Register("state", &State{})
}
