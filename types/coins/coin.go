package coin

import (
	"math/big"
	"kchain/types"
	"github.com/kooksee/kdb"
	"errors"
	"kchain/cmn"
	"encoding/json"
)

func NewCoin() *Coin {
	return &Coin{debug: false}
}

type Coin struct {
	To    string
	Coin  string
	debug bool
}

func (c *Coin) Decode(d []byte) error {
	return json.Unmarshal(d, c)
}

func (c *Coin) SetDebug(debug bool) *Coin {
	c.debug = debug
	return c
}

func (c *Coin) TransferFrom(from types.Address) error {
	to, err := cmn.ParsePubkey(c.To)
	if err != nil {
		return err
	}

	coin := big.NewFloat(0)
	cc, ok := coin.SetString(c.Coin)
	if !ok {
		return errors.New("coin parse error")
	}

	coin = cc

	return db.BatchUpdate(func(k *kdb.KHBatch) error {
		fAcc, err := k.Get(from.Bytes())
		if err != nil {
			return err
		}
		fCoin := big.NewFloat(0)
		fCoin, ok := fCoin.SetString(string(fAcc))
		if !ok {
			return errors.New(cmn.Fmt("coin %s parse error", fAcc))
		}

		tAcc, err := k.Get(to.Bytes())
		if err != nil {
			return err
		}
		tCoin := big.NewFloat(0)
		tCoin, ok = tCoin.SetString(string(tAcc))
		if !ok {
			return errors.New(cmn.Fmt("coin %s parse error", tAcc))
		}

		fCoin = fCoin.Sub(fCoin, coin)
		tCoin = tCoin.Add(tCoin, coin)

		if !c.debug {
			if err := k.Set(from.Bytes(), []byte(fCoin.String())); err != nil {
				return err
			}

			if err := k.Set(to.Bytes(), []byte(tCoin.String())); err != nil {
				return err
			}
		}

		return nil
	})
}
