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
	tx.Val.Power = t.Power
	tx.Val.PubKey = t.PubKey
	tx.Val.Address = t.Address
}

func UpdateValidator(val *types.Validator) {
	if val.Power <= 0 {
		if err := db.Del(val.Address); err != nil {
			panic(err.Error())
		}
	} else {
		dt, _ := val.Marshal()
		if err := db.Set(val.Address, dt); err != nil {
			panic(err.Error())
		}
	}
}

func GetValidator(addr []byte) (val *types.Validator) {
	dt, err := db.Get(addr)
	if err != nil {
		panic(err.Error())
	}
	if err := val.Unmarshal(dt); err != nil {
		panic(err.Error())
	}
	return
}

func Validators() (vals []*types.Validator) {
	if err := db.Range(func(_, value []byte) error {
		val := new(types.Validator)
		if err := val.Unmarshal(value); err != nil {
			panic(err.Error())
		}
		vals = append(vals, val)
		return nil
	}); err != nil {
		panic(err.Error())
	}
	return
}
