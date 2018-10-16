package wire

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	types.RegisterBlockAmino(cdc)
}

func GetCodec() *amino.Codec {
	return cdc
}

func Registry(fn func(cc *amino.Codec)) {
	fn(cdc)
}
