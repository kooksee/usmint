package kts

import (
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

type Validator struct {
	PubKey []byte `json:"pubkey,omitempty" mapstructure:"pubkey"`
	Power  int64  `json:"power,omitempty" mapstructure:"power"`
}

type IValidatorManager interface {
	// 更新验证者
	UpdateValidator(val *types.Validator) error

	// 检查验证者
	CheckValidator(val *types.Validator) error

	DecodeValidator([]byte) (*Validator, error)

	// 节点是否存在
	Has(address crypto.Address) bool
}
