package validator

import (
	"github.com/tendermint/tendermint/abci/types"
	"github.com/kooksee/usmint/kts"
	"fmt"
	"github.com/kooksee/usmint/mint/minter"
)

type Validator types.Validator

func (t *Validator) OnCheck(tx *kts.Transaction, res *types.ResponseCheckTx) {
	if t.Power <= 0 {
		t.Power = -1
	}

	if t.Power > 9 {
		res.Code = 1
		res.Log = fmt.Sprintf("max power is 10")
		return
	}

	if !minter.IsMaster(tx.GetMiner()) {
		res.Code = 1
		res.Log = fmt.Sprintf("the miner must be master")
		return
	}
}

func (t *Validator) OnDeliver(tx *kts.Transaction, res *types.ResponseDeliverTx) {
	tx.SetValidator(types.Validator{
		Power:   t.Power,
		PubKey:  t.PubKey,
		Address: t.Address,
	})
}
