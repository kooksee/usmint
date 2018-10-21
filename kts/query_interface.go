package kts

import (
	"github.com/kooksee/usmint/cmn/wire"
	"github.com/tendermint/tendermint/abci/types"
)

func DecodeQueryMsg(data []byte) (h QueryHandler, err error) {
	return h, wire.Decode(data, &h)
}

type QueryHandler interface {
	Do(res *types.ResponseQuery)
}
