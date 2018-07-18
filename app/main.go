package app

/*
实现tendermint abci的业务逻辑
 */

import (
	"github.com/tendermint/abci/types"
	"github.com/kooksee/kchain/types/code"
	"github.com/kooksee/usmint/mint"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/config"
)

type KApp struct {
	types.BaseApplication
	m *mint.Mint
}

func New() *KApp {
	// 要先初始化才能进行下一步操作
	mint.Init()
	return &KApp{m: mint.New()}
}

// 实现abci的Info协议
func (app *KApp) Info(req types.RequestInfo) (res types.ResponseInfo) {

	res.Data, _ = json.MarshalToString(map[string]interface{}{
		"name":     config.DefaultCfg().Moniker,
		"chain_id": "main",
	})
	res.LastBlockHeight = app.m.State().Height
	res.LastBlockAppHash = app.m.State().AppHash
	res.Version = req.Version

	return
}

// 实现abci的SetOption协议
func (app *KApp) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{Code: types.CodeTypeOK}
}

// 实现abci的DeliverTx协议
func (app *KApp) DeliverTx(txBytes []byte) types.ResponseDeliverTx {
	if err := app.m.DeliverTx(txBytes); err != nil {
		return types.ResponseDeliverTx{Code: code.ErrInternal.Code, Log: err.Error()}
	}
	return types.ResponseDeliverTx{Code: code.Ok.Code}
}

// 实现abci的CheckTx协议
func (app *KApp) CheckTx(txBytes []byte) types.ResponseCheckTx {
	if err := app.m.DeliverTx(txBytes); err != nil {
		return types.ResponseCheckTx{Code: code.ErrInternal.Code, Log: err.Error()}
	}
	return types.ResponseCheckTx{Code: code.Ok.Code}
}

// Commit will panic if InitChain was not called
func (app *KApp) Commit() types.ResponseCommit {
	return types.ResponseCommit{Data: app.m.Commit()}
}

func (app *KApp) Query(reqQuery types.RequestQuery) (res types.ResponseQuery) {
	return
}

// Save the validators in the merkle tree
func (app *KApp) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	cmn.MustNotErr("app InitChain error", app.m.InitChain(req.Validators...))
	return types.ResponseInitChain{}
}

func (app *KApp) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	cmn.ErrPipeLog("app BeginBlock error", app.m.BeginBlock(nil))
	return types.ResponseBeginBlock{}
}

func (app *KApp) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	val, err := app.m.EndBlock(nil)
	cmn.ErrPipeLog("app EndBlock error", err)
	return types.ResponseEndBlock{ValidatorUpdates: val}
}
