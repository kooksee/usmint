package app

/*
实现tendermint abci的业务逻辑
 */

import (
	"github.com/kooksee/usmint/mint"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/abci/types"
)

type KApp struct {
	types.BaseApplication
	m *mint.Mint
}

func New() *KApp {
	// 初始化mint模块
	mint.Init()
	return &KApp{m: mint.New()}
}

// 实现abci的Info协议
func (app *KApp) Info(req types.RequestInfo) (res types.ResponseInfo) {

	//res.Data = config.DefaultCfg().Moniker
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
	return app.m.DeliverTx(txBytes)
}

// 实现abci的CheckTx协议
func (app *KApp) CheckTx(txBytes []byte) types.ResponseCheckTx {
	return app.m.CheckTx(txBytes)
}

// Commit will panic if InitChain was not called
func (app *KApp) Commit() types.ResponseCommit {
	return types.ResponseCommit{Data: app.m.Commit()}
}

func (app *KApp) Query(reqQuery types.RequestQuery) (res types.ResponseQuery) {
	return app.m.QueryTx(reqQuery.Data)
}

// Save the validators in the merkle tree
func (app *KApp) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	cmn.MustNotErr(cmn.ErrPipe("app InitChain error", app.m.InitChain(req.Validators...)))
	return types.ResponseInitChain{}
}

func (app *KApp) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	cmn.ErrPipe("app BeginBlock error", app.m.BeginBlock(req))
	return types.ResponseBeginBlock{}
}

func (app *KApp) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	val, err := app.m.EndBlock(nil)
	cmn.ErrPipe("app EndBlock error", err)
	return types.ResponseEndBlock{ValidatorUpdates: val}
}
