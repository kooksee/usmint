package kts

import (
	"github.com/kooksee/usmint/cmn"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"github.com/tendermint/tendermint/abci/types"
	"crypto/ecdsa"
	"time"
	"github.com/kooksee/usmint/cmn/wire"
)

func NewTransaction() *Transaction {
	return &Transaction{Timestamp: uint64(time.Now().Unix())}
}

type Transaction struct {
	Sign      []byte `json:"sign"`
	NSign     []byte `json:"nsign"`
	Data      []byte `json:"data"`
	Event     string `json:"event"`
	Timestamp uint64 `json:"time"`
	miner     common.Address
	sender    common.Address
	val       types.Validator
}

func (t *Transaction) Decode(tx []byte) error {
	return wire.Decode(tx, t)
}

func (t *Transaction) SetValidator(val types.Validator) {
	t.val = val
}

func (t *Transaction) GetValidator() types.Validator {
	return t.val
}

func (t *Transaction) Encode() []byte {
	return wire.Encode(t)
}

func (t *Transaction) GetSender() common.Address {
	return t.sender
}

func (t *Transaction) GetSigHash() []byte {
	return crypto.Keccak256(t.Data, big.NewInt(int64(t.Timestamp)).Bytes())
}

func (t *Transaction) GetMiner() common.Address {
	return t.miner
}

func (t *Transaction) DoNSign(prv *ecdsa.PrivateKey) (err error) {
	t.NSign, err = crypto.Sign(t.GetSigHash(), prv)
	if err != nil {
		cmn.MustNotErr("DoNodeSign", err)
	}
	return
}

func (t *Transaction) DoSign(prv *ecdsa.PrivateKey) (err error) {
	t.Sign, err = crypto.Sign(crypto.Keccak256(t.Data), prv)
	if err != nil {
		cmn.MustNotErr("DoNodeSign", err)
	}
	return
}

// VerifySign 验证数据签名
func (t *Transaction) Verify() error {
	puk1, err := crypto.SigToPub(t.GetSigHash(), t.NSign)
	if err != nil {
		return cmn.ErrPipe("Transaction VerifySign Error With NSign", err)
	}
	t.miner = crypto.PubkeyToAddress(*puk1)

	puk, err := crypto.SigToPub(crypto.Keccak256(t.Data), t.Sign)
	if err != nil {
		return cmn.ErrPipe("Transaction VerifySign Error With Sign", err)
	}
	t.sender = crypto.PubkeyToAddress(*puk)

	return nil
}
