package tx

import (
	"strconv"
	"github.com/kooksee/kchain/types/cnst"
)


// Account
type Account struct {
	PubKey string `json:"pubkey,omitempty" mapstructure:"pubkey"`
	Power  int    `json:"power,omitempty" mapstructure:"power"`
	Ip     string `json:"power,omitempty" mapstructure:"ip"`
}

func NewAccount() *Account {
	return &Account{}
}

func (act *Account)ToSortBytes() []byte {
	return []byte(strconv.Itoa(act.Power) + act.PubKey)
}
func (act *Account)ToBytes() []byte {
	d, _ := json.Marshal(act)
	return d
}
func (act *Account) FromBytes(d []byte) error {
	return json.Unmarshal(d, act)
}
func (act *Account)GetPrefix() string {
	return cnst.AccountPrefix
}
func (act *Account) Key() []byte {
	return []byte(act.GetPrefix() + act.PubKey)
}
func (act *Account) ToKv() ([]byte, []byte) {
	return act.Key(), act.ToBytes()
}


