package minter

import (
	"github.com/tendermint/tendermint/libs/log"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/cmn/wire"
)

var (
	db     kdb.IKHash
	logger log.Logger
)

func init() {
	wire.Register("minter_set", &SetMiner{})
	wire.Register("minter_del", &DeleteMiner{})
}

func Init() {
	db = kdb.DefaultConfig().GetDb().KHash([]byte("minter"))
	logger = cmn.Log().With("pkg", "minter")
}
