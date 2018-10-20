package state

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn/consts"
	"github.com/kooksee/usmint/cmn/wire"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	db kdb.IKHash
)

func Init(logger log.Logger) {
	cfg := kdb.DefaultConfig()
	db = cfg.GetDb().KHash([]byte(consts.StatePrefix))
}

func init() {
	wire.Register("state", &State{})
}
