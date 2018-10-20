package app

/*
实现tendermint abci的业务逻辑
 */

import (
	"github.com/kooksee/usmint/mint"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/abci/types"
	"fmt"
	"github.com/kooksee/usmint/kts"
	"encoding/hex"
	"github.com/kooksee/usmint/mint/state"
	"github.com/kooksee/usmint/mint/minter"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/kooksee/usmint/usdb"
	"github.com/kooksee/usmint/node"
	"github.com/tendermint/tendermint/libs/log"
)

type KApp struct {
	valUpdates []types.Validator
	types.BaseApplication
	logger     log.Logger
}

func New(logger log.Logger) *KApp {
	// 初始化mint模块
	mint.Init(logger)
	return &KApp{logger: logger.With("module", "kapp")}
}

// 实现abci的Info协议
func (app *KApp) Info(req types.RequestInfo) (res types.ResponseInfo) {
	d, _ := cmn.JsonMarshalToString(req)
	app.logger.Info(d, "abci", "Info")

	res.Version = req.Version
	res.LastBlockHeight = state.GetState().Height
	res.LastBlockAppHash = state.GetState().AppHash

	return
}

// 实现abci的SetOption协议
func (app *KApp) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	d, _ := cmn.JsonMarshalToString(req)
	app.logger.Info(d, "abci", "SetOption")
	return types.ResponseSetOption{Code: types.CodeTypeOK}
}

// 实现abci的CheckTx协议
func (app *KApp) CheckTx(txBytes []byte) (res types.ResponseCheckTx) {
	app.logger.Info("abci.CheckTx", "tx", hex.EncodeToString(txBytes))

	// 获取tx
	txBytes = usdb.GetDb().Get(txBytes)

	// 检查tx大小
	if err := cmn.CheckMsgSize(txBytes); err != nil {
		res.Code = 1
		res.Log = err.Error()
		return
	}

	// 检查tx是否已经存在
	txHash := tmhash.Sum(txBytes)
	tx, err := node.GetNode().Indexer().Get(txHash)
	if err != nil {
		res.Code = 1
		res.Log = err.Error()
		return
	}

	if tx != nil && tx.Height != 0 && tx.Result.Code != 0 {
		res.Code = 1
		res.Log = fmt.Sprintf("the tx hash(%s) had existed", hex.EncodeToString(txHash))
		return
	}

	// decode tx
	tx1 := kts.NewTransaction()
	if err := tx1.Decode(txBytes); err != nil {
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
	if !minter.ExistMiner(tx1.GetMiner()) {
		res.Code = 1
		res.Log = fmt.Sprintf("the miner(%s) does not exist", tx1.GetMiner().Hash().String())
		return
	}

	// decode abci
	msg, err := kts.DecodeMsg(tx1.Data)
	if err != nil {
		res.Code = 1
		res.Log = fmt.Sprintf("Mint CheckTx decodeMsg Error(%s)", err.Error())
		return
	}

	msg.OnCheck(tx1, &res)

	return
}

// 实现abci的DeliverTx协议
func (app *KApp) DeliverTx(txBytes []byte) (res types.ResponseDeliverTx) {
	app.logger.Info("abci.DeliverTx", "tx", hex.EncodeToString(txBytes))

	// 获取tx
	txBytes = usdb.GetDb().Get(txBytes)

	tx := kts.NewTransaction()
	tx.Decode(txBytes)

	msg, _ := kts.DecodeMsg(tx.Data)
	msg.OnDeliver(tx, &res)

	if res.Code == 0 {
		if tx.Event == "validator" {
			app.valUpdates = append(app.valUpdates, tx.GetValidator())
		}
		state.GetState().AppHash = cmn.Ripemd160(append(state.GetState().AppHash, txBytes...))
	}
	return
}

// Commit will panic if InitChain was not called
func (app *KApp) Commit() types.ResponseCommit {
	app.logger.Info("abci.Commit", "apphash", hex.EncodeToString(state.GetState().AppHash))

	return types.ResponseCommit{Data: state.GetState().AppHash}
}

func (app *KApp) Query(reqQuery types.RequestQuery) (res types.ResponseQuery) {
	app.logger.Info("abci.Query")

	// 检查tx大小
	if err := cmn.CheckMsgSize(reqQuery.Data); err != nil {
		res.Code = 1
		res.Log = err.Error()
		return
	}

	msg, err := kts.DecodeQueryMsg(reqQuery.Data)
	if err != nil {
		res.Code = 1
		res.Log = err.Error()
		return
	}
	msg.Do(&res)
	return
}

// Save the validators in the merkle tree
func (app *KApp) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	d, _ := cmn.JsonMarshalToString(req)
	app.logger.Info(d, "abci", "InitChain")

	// 添加master miner address
	minter.InitMaster()
	return types.ResponseInitChain{}
}

func (app *KApp) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	d, _ := cmn.JsonMarshalToString(req)
	app.logger.Info(d, "abci", "BeginBlock")

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
	d, _ := cmn.JsonMarshalToString(req)
	app.logger.Info(d, "abci", "EndBlock")

	return types.ResponseEndBlock{ValidatorUpdates: app.valUpdates}
}
