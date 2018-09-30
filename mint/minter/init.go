package minter

import (
	"github.com/tendermint/tendermint/libs/log"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/wire"
)

var (
	db     kdb.IKHash
	logger log.Logger
)

func init() {
	cc := wire.GetCodec()
	cc.RegisterConcrete(&SetMiner{}, "mint/minter/set", nil)
	cc.RegisterConcrete(&DeleteMiner{}, "mint/minter/delete", nil)
}

func Init() {
	db = kdb.DefaultConfig().GetDb().KHash([]byte("minter"))
	logger = cmn.Log().With("pkg", "minter")
}
