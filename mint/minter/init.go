package minter

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn/wire"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/ethereum/go-ethereum/common"
)

var (
	db kdb.IKHash
)

func init() {
	wire.Register("miner", &Miner{})
	wire.Register("miner_set", &SetMiner{})
	wire.Register("miner_del", &DeleteMiner{})
	wire.Register("address", &common.Address{})
}

func Init(logger log.Logger) {
	db = kdb.DefaultConfig().GetDb().KHash([]byte("miner"))
}
