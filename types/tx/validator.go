package tx

import (
	"strconv"
	"github.com/kooksee/kchain/types/cnst"
)


// Validator
type Validator struct {
	PubKey string `json:"pubkey,omitempty" mapstructure:"pubkey"`
	Power  int    `json:"power,omitempty" mapstructure:"power"`
}

func NewValidator() *Validator {
	return &Validator{}
}

func (val *Validator)ToSortBytes() []byte {
	return []byte(strconv.Itoa(val.Power) + val.PubKey)
}
func (val *Validator)ToBytes() []byte {
	d, _ := json.Marshal(val)
	return d
}
func (val *Validator) FromBytes(d []byte) error {
	return json.Unmarshal(d, val)
}
func (val *Validator)GetPrefix() string {
	return cnst.ValidatorPrefix
}

func (val *Validator)Key() []byte {
	return []byte(val.GetPrefix() + val.PubKey)
}

func (val *Validator) ToKv() ([]byte, []byte) {
	return val.Key(), val.ToBytes()
}
