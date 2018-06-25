package types

import (
	"encoding/hex"
	"github.com/tendermint/go-crypto"
	"fmt"
	"errors"
	"github.com/kooksee/usmint/cmn"
	"time"
)

func NewTransaction() *Transaction {
	return &Transaction{}
}

type Transaction struct {
	PubKey    string `json:"pubkey,omitempty"`
	Signature string `json:"sign,omitempty"`
	Data      string `json:"data,omitempty"`
	Event     string `json:"event,omitempty"`

	timestamp int64
	hash      []byte
	pubkey    crypto.PubKey
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

	return []byte(fmt.Sprintf("%s%s%d", t.Data, t.Event, time.Now().Second()))
}

// Sign 签名
func (t *Transaction) Sign(priv crypto.PrivKey) ([]byte, error) {
	msg := t.signMsg()
	if msg == nil {
		return nil, errors.New("签名数据为空")
	}

	sign := priv.Sign(msg)
	if sign.IsZero() {
		return nil, errors.New("签名失败")
	}

	return sign.Bytes(), nil
}

func (t *Transaction) GetPubKey() (crypto.PubKey, error) {
	if t.pubkey != nil {
		return t.pubkey, nil
	}
	return cmn.ParsePubkey(t.PubKey)
}

// VerifySign 签名验证
func (t *Transaction) VerifySign() error {

	if t.Signature == "" || t.PubKey == "" {
		return errors.New("sign or pubkey is null")
	}

	// 区块签名验证
	d, err := hex.DecodeString(t.PubKey)
	if err != nil {
		return err
	}

	pk, err := crypto.PubKeyFromBytes(d)
	if err != nil {
		return err
	}

	// 缓存pubkey
	t.pubkey = pk

	sign, err := hex.DecodeString(t.Signature)
	if err != nil {
		return err
	}

	sig, err := crypto.SignatureFromBytes(sign)
	if err != nil {
		return err
	}

	msg := t.signMsg()
	if msg == nil {
		return errors.New("签名数据为空")
	}

	if !pk.VerifyBytes(crypto.Ripemd160(msg), sig) {
		return errors.New("transaction verify false")
	}

	return nil
}
