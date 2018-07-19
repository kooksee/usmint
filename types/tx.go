package types

import (
	"encoding/hex"
	"fmt"
	"errors"
	"github.com/tendermint/go-crypto"
)

func DecodeTx(bs []byte) (*Transaction, error) {
	tx := NewTransaction()
	return tx, json.Unmarshal(bs, tx)
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

type Transaction struct {
	Signature     string `json:"sign,omitempty"`
	NodeSignature string `json:"node_sign,omitempty"`
	Pubkey        string `json:"pubkey,omitempty"`
	Data          string `json:"data,omitempty"`
	Event         string `json:"event,omitempty"`
	Timestamp     uint64 `json:"time,omitempty"`
	hash          []byte
	pubkey        crypto.PubKey
}

func (t *Transaction) Dumps() ([]byte, error) {
	return json.Marshal(t)
}

// FromBytes 解析Transaction
func (t *Transaction) Decode(bs []byte) error {
	return json.Unmarshal(bs, t)
}

// Hash tx hash
func (t *Transaction) Hash() ([]byte, error) {
	if len(t.hash) != 0 {
		return t.hash, nil
	}

	if len(t.Signature) == 0 {
		return nil, errors.New("签名为空")
	}

	sign, err := hex.DecodeString(t.Signature)
	if err != nil {
		return nil, err
	}

	t.hash = crypto.Ripemd160(sign)
	return t.hash, nil
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

// Sign 签名
func (t *Transaction) Sign(priv crypto.PrivKey) ([]byte, error) {
	msg := t.signMsg()
	if msg == nil {
		return nil, errors.New("签名数据为空")
	}

	return priv.Sign(msg).Bytes(), nil
}

func (t *Transaction) GetPubKey() (crypto.PubKey, error) {
	if t.pubkey != nil {
		return t.pubkey, nil
	}

	pubkey, err := hex.DecodeString(t.Pubkey)
	if err != nil {
		return nil, err
	}

	pk, err := crypto.PubKeyFromBytes(pubkey)
	if err != nil {
		return nil, err
	}

	t.pubkey = pk
	return pk, nil
}

// VerifySign 签名验证
func (t *Transaction) VerifySign() error {
	sign, err := hex.DecodeString(t.Signature)
	if err != nil {
		return err
	}

	s, err := crypto.SignatureFromBytes(sign)
	if err != nil {
		return err
	}

	pubkey, err := hex.DecodeString(t.Pubkey)
	if err != nil {
		return err
	}

	pk, err := crypto.PubKeyFromBytes(pubkey)
	if err != nil {
		return err
	}

	if !pk.VerifyBytes(t.signMsg(), s) {
		return errors.New("transaction verify false")
	}

	t.pubkey = pk

	return nil
}
