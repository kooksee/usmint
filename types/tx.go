package types

import (
	"encoding/hex"
	"github.com/tendermint/go-crypto"
)

type Transaction struct {
	PubKey    string `json:"pubKey,omitempty" validate:"required"`
	Signature string `json:"sign,omitempty"   validate:"required"`
	Method    string `json:"method,omitempty" validate:"required"`
	TxID      string `json:"txId,omitempty"`
	Data      string `json:"data,omitempty"`
	Address   string `json:"addr,omitempty"`
	Timestamp int64  `json:"time,omitempty"`
}

func (t *Transaction) Dumps() []byte {
	d, _ := json.Marshal(t)
	return d
}

func (t *Transaction) GetTxID() string {
	d, _ := hex.DecodeString(t.Signature)
	return hex.EncodeToString(crypto.Ripemd160(d))
}

func NewTransaction() *Transaction {
	return &Transaction{}
}
