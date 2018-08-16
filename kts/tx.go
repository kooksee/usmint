package kts

import (
	"fmt"
	"errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/kooksee/usmint/cmn"
	"encoding/hex"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
)

func DecodeTx(bs []byte) (*Transaction, error) {
	tx := NewTransaction()
	return tx, cmn.ErrPipe("DecodeTx2", cmn.JsonUnmarshal(bs, tx))
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
	pubkey        crypto.PubKey
}

func (t *Transaction) Dumps() ([]byte, error) {
	return cmn.JsonMarshal(t)
}

// FromBytes 解析Transaction
func (t *Transaction) Decode(bs []byte) error {
	return cmn.JsonUnmarshal(bs, t)
}

func (t *Transaction) SignMsg() []byte {
	if t.Data == "" {
		return nil
	}

	if t.Event == "" {
		return nil
	}

	if t.Timestamp == 0 {
		return nil
	}

	return crypto.Sha256([]byte(fmt.Sprintf("%s%s%d", t.Data, t.Event, t.Timestamp)))
}

func (t *Transaction) GetPubkey() crypto.PubKey {
	return t.pubkey
}

// VerifySign 验证数据签名
func (t *Transaction) VerifySign() (addr common.Address, err error) {
	sign, err := hex.DecodeString(t.Signature)
	if err != nil {
		return addr, cmn.ErrPipe("Transaction VerifyNodeSign 1", err)
	}

	data, err := hex.DecodeString(t.Data)
	if err != nil {
		return addr, cmn.ErrPipe("Transaction VerifyNodeSign 2", err)
	}

	pubk, err := ecrypto.SigToPub(data, sign)
	if err != nil {
		return addr, cmn.ErrPipe("Transaction VerifyNodeSign 3", err)
	}

	return ecrypto.PubkeyToAddress(*pubk), nil
}

// VerifyNodeSign 节点签名验证
func (t *Transaction) VerifyNodeSign() error {
	sign, err := hex.DecodeString(t.NodeSignature)
	if err != nil {
		return cmn.ErrPipe("Transaction VerifyNodeSign 1", err)
	}

	pubkey, err := hex.DecodeString(t.Pubkey)
	if err != nil {
		return cmn.ErrPipe("Transaction VerifyNodeSign 2", err)
	}

	pk, err := cryptoAmino.PubKeyFromBytes(pubkey)
	if err != nil {
		return cmn.ErrPipe("Transaction VerifyNodeSign PubKeyFromBytes", err)
	}

	if !pk.VerifyBytes(t.SignMsg(), sign) {
		return cmn.ErrPipe("Transaction VerifyNodeSign 4", errors.New("transaction verify false"))
	}

	t.pubkey = pk

	return nil
}
