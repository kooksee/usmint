package types

import (
	"encoding/hex"
	"fmt"
	"errors"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
)

func DecodeTx(bs []byte) (*Transaction, error) {
	tx := &Transaction{}
	return tx, json.Unmarshal(bs, tx)
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

type Transaction struct {
	Signature string `json:"sign,omitempty"`
	Data      string `json:"data,omitempty"`
	Event     string `json:"event,omitempty"`
	Timestamp uint64 `json:"time,omitempty"`
	hash      []byte
	pubkey    *ecdsa.PublicKey
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

	t.hash = crypto.Keccak256(sign)

	return t.hash, nil
}

func (t *Transaction) signMsg() []byte {
	if len(t.Data) == 0 {
		return nil
	}

	if len(t.Event) == 0 {
		return nil
	}

	return []byte(fmt.Sprintf("%s%s%d", t.Data, t.Event, t.Timestamp))
}

// Sign 签名
func (t *Transaction) Sign(priv *ecdsa.PrivateKey) ([]byte, error) {
	msg := t.signMsg()
	if msg == nil {
		return nil, errors.New("签名数据为空")
	}

	sig, err := crypto.Sign(msg, priv)
	if err != nil {
		return nil, errors.New("签名失败")
	}

	return sig, nil
}

func (t *Transaction) GetPubKey() *ecdsa.PublicKey {
	return t.pubkey
}

// VerifySign 签名验证
func (t *Transaction) VerifySign() error {
	sign, err := hex.DecodeString(t.Signature)
	if err != nil {
		return err
	}

	pubkey, err := crypto.Ecrecover(t.signMsg(), sign)
	if err != nil {
		return errors.New("transaction verify false")
	}
	t.pubkey = crypto.ToECDSAPub(pubkey)

	return nil
}
