package kts

import (
	"fmt"
	"errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func DecodeTx(bs []byte) (*Transaction, error) {
	tx := NewTransaction()
	return tx, cmn.JsonUnmarshal(bs, tx)
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

type Transaction struct {
	Signature     string `json:"sign,omitempty"`
	NodeSignature string `json:"node_sign,omitempty"`
	Pubkey        string `json:"pubkey,omitempty"`
	Data          []byte `json:"data,omitempty"`
	Event         string `json:"event,omitempty"`
	Timestamp     uint64 `json:"time,omitempty"`
	pubkey        crypto.PubKey
}

func (t *Transaction) Dumps() ([]byte, error) {
	return cmn.JsonMarshal(t)
}

// FromBytes 解析Transaction
func (t *Transaction) Decode(bs []byte) error {
	return cmn.JsonUnmarshal(bs, t)
}

func (t *Transaction) signMsg() []byte {
	if len(t.Data) == 0 {
		return nil
	}

	if len(t.Event) == 0 {
		return nil
	}

	if t.Timestamp == 0 {
		return nil
	}

	return crypto.Ripemd160([]byte(fmt.Sprintf("%s%s%d", t.Data, t.Event, t.Timestamp)))
}

func (t *Transaction) GetPubkey() crypto.PubKey {
	return t.pubkey
}

// VerifySign 签名验证
func (t *Transaction) VerifySign() error {
	sign, err := hexutil.Decode(t.Signature)
	if err != nil {
		return cmn.ErrPipe("Transaction.VerifySign.hexutil.Decode error", err)
	}

	s, err := cryptoAmino.SignatureFromBytes(sign)
	if err != nil {
		return cmn.ErrPipe("Transaction.VerifySign.cryptoAmino.SignatureFromBytes error", err)
	}

	pubkey, err := hexutil.Decode(t.Pubkey)
	if err != nil {
		return cmn.ErrPipe("Transaction.VerifySign.Decode.Pubkey error", err)
	}

	pk, err := cryptoAmino.PubKeyFromBytes(pubkey)
	if err != nil {
		return cmn.ErrPipe("Transaction.VerifySign.cryptoAmino.PubKeyFromBytes error", err)
	}

	if !pk.VerifyBytes(t.signMsg(), s) {
		return cmn.ErrPipe("Transaction.VerifySign", errors.New("transaction verify false"))
	}

	t.pubkey = pk

	return nil
}
