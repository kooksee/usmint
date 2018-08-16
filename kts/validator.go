package kts

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/kooksee/usmint/cmn"
	"encoding/hex"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

func PubkeyFrom(data []byte) (crypto.PubKey, error) {
	pk, err := cryptoAmino.PubKeyFromBytes(data)
	return pk, cmn.ErrPipe("kts.PubkeyFrom", err)
}

func DecodeValidator(data string) (*Validator, error) {
	dt, err := hex.DecodeString(data)
	if err != nil {
		return nil, cmn.ErrPipe("DecodeValidator 1", err)
	}

	val := &Validator{}
	if err := cmn.ErrPipe("DecodeValidator 2", cmn.JsonUnmarshal(dt, val)); err != nil {
		return nil, err
	}

	val.address, err = hex.DecodeString(val.Address)
	return val, cmn.ErrPipe("DecodeValidator 3", err)
}

type Validator struct {
	Name       string       `json:"name,omitempty"`
	Address    string       `json:"address,omitempty"`
	PubKey     types.PubKey `json:"pub_key,omitempty"`
	Power      int64        `json:"power,omitempty"`
	CreateTime int64        `json:"create_time,omitempty"`
	UpdateTime int64        `json:"update_time,omitempty"`
	address    []byte
}

func (v *Validator) Encode() ([]byte, error) {
	return cmn.JsonMarshal(v)
}

func (v *Validator) GetAddress() []byte {
	return v.address
}

type IValidatorManager interface {
	// 更新验证者
	UpdateValidator(val *Validator) error

	// 检查验证者
	CheckValidator(val *Validator) error

	DecodeValidator([]byte) (*Validator, error)

	// 节点是否存在
	Has(address crypto.Address) bool
}
