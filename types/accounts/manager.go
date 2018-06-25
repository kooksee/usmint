package accounts

import (
	"kchain/types"
	"errors"
)

func GetAccount(addr types.Address) (*Account, error) {
	k := Prefix + addr.String()
	val, ok := cache.Get(k)
	if ok {
		return val.(*Account), nil
	}

	ok, err := db.Exist(addr.Bytes())
	if err != nil {
		return nil, err
	}

	if ok {
		val, err := db.Get(addr.Bytes())
		if err != nil {
			return nil, err
		}

		acc := NewAccount()
		if err := acc.Decode(val); err != nil {
			return nil, err
		}
		cache.SetDefault(k, acc)
		return acc, nil
	}

	return nil, errors.New("找不到该account")
}
