package kts

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/kooksee/usmint/cmn"
	"encoding/hex"
)

func DecodeValidator(data string) (*Validator, error) {
	dt, err := hex.DecodeString(data)
	if err != nil {
		return nil, cmn.ErrPipe("DecodeValidator 1", err)
	}

	val := &Validator{}
	return val, cmn.ErrPipe("DecodeValidator 2", cmn.JsonUnmarshal(dt, val))
}

type Validator struct {
	Name       string `json:"name,omitempty"`
	Address    string `json:"address,omitempty"`
	Power      int64  `json:"power,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
	UpdateTime int64  `json:"update_time,omitempty"`
}

func (v *Validator) Encode() ([]byte, error) {
	return cmn.JsonMarshal(v)
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
