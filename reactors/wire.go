package reactors

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterRTx(cdc)
	types.RegisterBlockAmino(cdc)
}
