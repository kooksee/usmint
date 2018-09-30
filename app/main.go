package app

/*
实现tendermint abci的业务逻辑
 */

import (
	"github.com/kooksee/usmint/mint"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/abci/types"
	tt "github.com/tendermint/tendermint/types"
	"fmt"
	"github.com/kooksee/usmint/kts"
	"github.com/kooksee/usmint/wire"
	"encoding/hex"
	"github.com/kooksee/usmint/mint/validator"
	"github.com/kooksee/usmint/mint/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kooksee/usmint/mint/minter"
)

func checkMsgSize(txBytes []byte) error {
	if len(txBytes) > maxMsgSize {
		return fmt.Errorf("msg size exceeds max size (%d > %d)", len(txBytes), maxMsgSize)
	}
	return nil
}

var maxMsgSize = 1024 * 1024

func decodeMsg(bz []byte) (msg kts.DataHandler, err error) {
	return msg, wire.GetCodec().UnmarshalBinaryBare(bz, &msg)
}

type KApp struct {
	types.BaseApplication

	valUpdates []types.Validator
}

func New() *KApp {
	// 初始化mint模块

	mint.Init()
	return &KApp{}
}

// 实现abci的Info协议
func (app *KApp) Info(req types.RequestInfo) (res types.ResponseInfo) {

	res.Version = req.Version
	res.LastBlockHeight = state.GetState().Height
	res.LastBlockAppHash = state.GetState().AppHash
	//res.Data, _ = cmn.JsonMarshalToString(&map[string]interface{}{
	//	"config":    cmn.GetCfg(),
	//	"node_info": cmn.GetNode().NodeInfo(),
	//})

	return
}

// 实现abci的SetOption协议
func (app *KApp) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{Code: types.CodeTypeOK}
}

// 实现abci的DeliverTx协议
func (app *KApp) DeliverTx(txBytes []byte) (res types.ResponseDeliverTx) {
	tx := kts.NewTransaction()
	if err := wire.GetCodec().UnmarshalBinaryBare(txBytes, tx); err != nil {
		res.Code = 1
		res.Log = fmt.Sprintf("tx decode error(%s)", err.Error())
		return
	}

	if tx.Event == "valUpdate" {
		val := &types.Validator{}
		if err := val.Unmarshal(tx.Data); err != nil {
			return
		}

		if int(val.Power) > 9 {
			return
		}

		validator.UpdateValidator(val)
		app.valUpdates = append(app.valUpdates, *val)
	}

	msg, err := decodeMsg(tx.Data)
	if err != nil {
		res.Code = 1
		res.Log = cmn.ErrPipe("Mint DeliverTx decodeMsg Error", err).Error()
		return
	}

	msg.OnDeliver(tx, &res)
	state.GetState().AppHash = crypto.Keccak256(state.GetState().AppHash, txBytes)
	return
}

// 实现abci的CheckTx协议
func (app *KApp) CheckTx(txBytes []byte) (res types.ResponseCheckTx) {
	// 检查tx大小
	if err := checkMsgSize(txBytes); err != nil {
		res.Code = 1
		res.Log = err.Error()
		return
	}

	// 检查tx是否已经存在
	txHash := tt.Tx(txBytes).Hash()
	tx, _ := cmn.GetNode().Indexer().Get(txHash)
	if tx != nil {
		res.Code = 1
		res.Log = fmt.Sprintf("the hash(%s) had existed", hex.EncodeToString(txHash))
		return
	}

	// decode tx
	tx1 := kts.NewTransaction()
	if err := wire.GetCodec().UnmarshalBinaryBare(txBytes, tx1); err != nil {
		res.Code = 1
		res.Log = fmt.Sprintf("tx decode error(%s)", err.Error())
		return
	}

	// tx verify
	if err := tx1.Verify(); err != nil {
		res.Code = 1
		res.Log = cmn.ErrPipe("Mint CheckTx VerifySign Error", err).Error()
		return
	}

	// check miner auth
	if !minter.Exist(tx1.GetMiner()) {
		res.Code = 1
		res.Log = fmt.Sprintf("the miner(%s) does not exist", tx1.GetMiner().Hash().String())
		return
	}

	// decode abci
	msg, err := decodeMsg(tx1.Data)
	if err != nil {
		res.Code = 1
		res.Log = cmn.ErrPipe("Mint CheckTx decodeMsg Error", err).Error()
		return
	}

	msg.OnCheck(tx1, &res)
	return
}

// Commit will panic if InitChain was not called
func (app *KApp) Commit() types.ResponseCommit {
	return types.ResponseCommit{Data: state.GetState().AppHash}
}

func (app *KApp) Query(reqQuery types.RequestQuery) (res types.ResponseQuery) {
	// 检查tx大小
	if err := checkMsgSize(reqQuery.Data); err != nil {
		res.Code = 1
		res.Log = err.Error()
		return
	}

	msg, err := decodeMsg(reqQuery.Data)
	if err != nil {
		res.Code = 1
		res.Log = cmn.ErrPipe("Mint CheckTx VerifySign Error", err).Error()
		return
	}

	msg.OnQuery(&res)
	return
}

// Save the validators in the merkle tree
func (app *KApp) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	for _, val := range req.Validators {
		validator.UpdateValidator(&val)
	}

	// 添加master
	minter.InitMaster()
	return types.ResponseInitChain{}
}

func (app *KApp) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	app.valUpdates = make([]types.Validator, 0)

	st := state.GetState()
	st.Height = req.Header.Height
	st.BlockHash = req.Hash
	st.AppHash = req.Header.AppHash
	st.Time = req.Header.Time
	st.NumTxs = req.Header.NumTxs
	st.TotalTxs = req.Header.TotalTxs

	return types.ResponseBeginBlock{}
}

func (app *KApp) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{ValidatorUpdates: app.valUpdates}
}
