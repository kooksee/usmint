package validator

import (
	"github.com/tendermint/tendermint/abci/types"
)

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
