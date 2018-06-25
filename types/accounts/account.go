package accounts

import (
	"encoding/hex"
	"kchain/cmn"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/common"
	"kchain/types"
)

// Account
type Account struct {
	PubKey string          `json:"pubkey,omitempty" mapstructure:"pubkey"`
	Roles  map[string]bool `json:"roles,omitempty" mapstructure:"roles"`

	pubKey  crypto.PubKey
	address common.HexBytes
}

func NewAccount() *Account {
	return &Account{}
}

func (act *Account) ToSortBytes() ([]byte, error) {
	return cmn.SortStruct(act)
}

func (a *Account) Encode() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Account) Decode(d []byte) error {
	return json.Unmarshal(d, a)
}

func (act *Account) Key() ([]byte, error) {
	addr, err := act.GetAddress()
	if err != nil {
		return nil, err
	}
	return append([]byte(Prefix), addr.Bytes()...), nil
}

func (a *Account) Save() error {
	val, err := a.Encode()
	if err != nil {
		return err
	}

	addr, err := a.GetAddress()
	if err != nil {
		return err
	}

	return db.Set(addr.Bytes(), val)
}

// create an Address from a string
func (a *Account) GetAddress() (addr types.Address, err error) {
	if len(a.address) == 0 {
		pk, err := a.GetPubKey()
		if err != nil {
			return nil, err
		}

		a.address = pk.Address()
	}
	return a.address, nil
}

func (a *Account) GetPubKey() (crypto.PubKey, error) {
	if a.pubKey == nil {
		d, err := hex.DecodeString(a.PubKey)
		if err != nil {
			return nil, err
		}

		pk, err := crypto.PubKeyFromBytes(d)
		if err != nil {
			return nil, err
		}

		a.pubKey = pk
	}

	return a.pubKey, nil
}
