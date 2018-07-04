package mint

import (
	"encoding/hex"
	"crypto/ecdsa"
	"ybkchain/cmn"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Account
type Account struct {
	PubKey string          `json:"pubkey,omitempty" mapstructure:"pubkey"`
	Roles  map[string]bool `json:"roles,omitempty" mapstructure:"roles"`

	pubKey  *ecdsa.PublicKey
	address common.Address
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
func (a *Account) GetAddress() (addr common.Address, err error) {
	if len(a.address) == 0 {
		a.address = common.HexToAddress(a.PubKey)
	}
	return a.address, nil
}

func (a *Account) GetPubKey() (*ecdsa.PublicKey, error) {
	if a.pubKey == nil {
		d, err := hex.DecodeString(a.PubKey)
		if err != nil {
			return nil, err
		}
		a.pubKey = crypto.ToECDSAPub(d)
	}

	return a.pubKey, nil
}
