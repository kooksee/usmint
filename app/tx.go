package app

import (
	"github.com/tendermint/go-crypto"

	"github.com/pkg/errors"
	"encoding/hex"
	"github.com/kooksee/kchain/types/cnst"
	"github.com/tendermint/tmlibs/common"
	"strings"
	"github.com/kooksee/kchain/types"
)

type Transaction types.Transaction

func NewTransaction() *Transaction {
	return &Transaction{}
}

// FromBytes 解析Transaction
func (t *Transaction) FromBytes(bs []byte) error {
	return json.Unmarshal(bs, t)
}

func (t *Transaction) GetTxID() string {
	d, _ := hex.DecodeString(t.Signature)
	return hex.EncodeToString(crypto.Ripemd160(d))
}

// ToBytes Marshal
func (t *Transaction) ToBytes() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Transaction) IsRpcOpt() bool {
	return strings.HasPrefix(t.Method, "rpc")
}

func (t *Transaction) IsCTTOpt() bool {
	return strings.HasPrefix(t.Method, "ctt")
}

func (t *Transaction) Verify() error {

	if t.Signature == "" || t.PubKey == "" {
		return errors.New("sign or pubkey is null")
	}

	// 检查发送tx的节点有没有在区块链中
	if !state.db.Has([]byte(cnst.ValidatorPrefix + t.PubKey)) {
		return errors.New(common.Fmt("the node %s does not exist", t.PubKey))
	}

	// 区块签名验证
	d, _ := hex.DecodeString(t.PubKey)
	if pk, err := crypto.PubKeyFromBytes(d); err != nil {
		return err
	} else {
		d, _ := hex.DecodeString(t.Signature)
		if sig, err := crypto.SignatureFromBytes(d); err != nil {
			return err
		} else {
			if !pk.VerifyBytes(crypto.Ripemd160([]byte(t.Data)), sig) {
				return errors.New("transaction verify false")
			}
		}
	}
	return nil
}
