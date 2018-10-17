package kts

import (
	"github.com/kooksee/usmint/cmn/wire"
	"github.com/tendermint/tendermint/abci/types"
)

func DecodeMsg(data []byte) (h DataHandler, err error) {
	return h, wire.Decode(data, h)
}

type DataHandler interface {
	OnCheck(tx *Transaction, res *types.ResponseCheckTx)
	OnDeliver(tx *Transaction, res *types.ResponseDeliverTx)
}

type BaseDataHandler struct {
	DataHandler
}

func (t *BaseDataHandler) OnCheck(tx *Transaction, res *types.ResponseCheckTx)     {}
func (t *BaseDataHandler) OnDeliver(tx *Transaction, res *types.ResponseDeliverTx) {}
