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

func DecodeQueryMsg(data []byte) (h QueryHandler, err error) {
	return h, wire.Decode(data, h)
}

type QueryHandler interface {
	Do(res *types.ResponseQuery)
}

func DecodeMsg(data []byte) (h DataHandler, err error) {
	return h, wire.Decode(data, h)
}

type DataHandler interface {
	OnCheck(tx *Transaction, res *types.ResponseCheckTx)
	OnDeliver(tx *Transaction, res *types.ResponseDeliverTx)
}

type BaseDataHandler struct {
	DataHandler
}

func (t *BaseDataHandler) OnCheck(tx *Transaction, res *types.ResponseCheckTx)     {}
func (t *BaseDataHandler) OnDeliver(tx *Transaction, res *types.ResponseDeliverTx) {}

func NewTransaction() *Transaction {
	return &Transaction{Timestamp: uint64(time.Now().Unix())}
}

type Transaction struct {
	Sign      []byte          `json:"sign"`
	NodeSign  []byte          `json:"node_sign"`
	Data      []byte          `json:"data"`
	Event     string          `json:"event"`
	Timestamp uint64          `json:"time"`
	miner     common.Address
	sender    common.Address
	Val       types.Validator `json:"-"`
}

func (t *Transaction) Decode(tx []byte) error {
	return wire.Decode(tx, t)
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

func (t *Transaction) DoNodeSign(prv *ecdsa.PrivateKey) (err error) {
	t.NodeSign, err = crypto.Sign(t.GetSigHash(), prv)
	if err != nil {
		cmn.MustNotErr("DoNodeSign", err)
	}
	return
}

func (t *Transaction) DoSenderSign(prv *ecdsa.PrivateKey) (err error) {
	t.Sign, err = crypto.Sign(crypto.Keccak256(t.Data), prv)
	if err != nil {
		cmn.MustNotErr("DoNodeSign", err)
	}
	return
}

// VerifySign 验证数据签名
func (t *Transaction) Verify() error {
	puk1, err := crypto.SigToPub(t.GetSigHash(), t.NodeSign)
	if err != nil {
		return cmn.ErrPipe("Transaction VerifySign error with node", err)
	}
	t.miner = crypto.PubkeyToAddress(*puk1)

	puk, err := crypto.SigToPub(crypto.Keccak256(t.Data), t.Sign)
	if err != nil {
		return cmn.ErrPipe("Transaction VerifySign error with data", err)
	}
	t.sender = crypto.PubkeyToAddress(*puk)

	return nil
}

func init() {
	wire.RegisterInterface((*DataHandler)(nil))
	wire.Register("tx", &Transaction{})
	wire.Register("baseHandler", &BaseDataHandler{})
}
