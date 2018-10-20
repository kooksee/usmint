package minter

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn/wire"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	db kdb.IKHash
)

func init() {
	wire.Register("miner_set", &SetMiner{})
	wire.Register("miner_del", &DeleteMiner{})
}

func Init(logger log.Logger) {
	db = kdb.DefaultConfig().GetDb().KHash([]byte("miner"))
}
