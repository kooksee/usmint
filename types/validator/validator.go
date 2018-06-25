package validator

import (
	"github.com/tendermint/go-crypto"
	"encoding/hex"
	"errors"
)

func NewValidator() *Validator {
	return &Validator{}
}

// Validator
type Validator struct {
	PubKey string `json:"pubkey,omitempty" mapstructure:"pubkey"`
	Power  int64  `json:"power,omitempty" mapstructure:"power"`
	pubkey crypto.PubKey
}

func (v *Validator) GetPubkey() (crypto.PubKey, error) {
	if v.pubkey != nil {
		return v.pubkey, nil
	}

	d, err := hex.DecodeString(v.PubKey)
	if err != nil {
		return nil, err
	}

	pk, err := crypto.PubKeyFromBytes(d)
	if err != nil {
		return nil, err
	}

	v.pubkey = pk
	return pk, nil
}

func (v *Validator) Has() (bool, error) {
	pk, err := v.GetPubkey()
	if err != nil {
		return false, err
	}
	return db.Exist(pk.Address().Bytes())
}

// Check 检查Power值和Pubkey
func (v *Validator) Check() error {
	if v.Power > 9 {
		return errors.New("Power值超过限制")
	}

	_, err := v.GetPubkey()
	if err != nil {
		return err
	}
	return nil
}

func (v *Validator) Decode(val []byte) error {
	return json.Unmarshal(val, v)
}

func (v *Validator) Delete() error {
	pk, err := v.GetPubkey()
	if err != nil {
		return err
	}
	return db.Del(pk.Address().Bytes())
}

func (v *Validator) Save() error {
	pk, err := v.GetPubkey()
	if err != nil {
		return err
	}

	val, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return db.Set(pk.Address().Bytes(), val)
}
