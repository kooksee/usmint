package app

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/kooksee/kchain/types/cnst"
	"github.com/kooksee/kchain/types/code"
)

type KApp struct {
	types.BaseApplication
	ValUpdates []types.Validator
	GValidator string
}

func New(name, dbDir string) *KApp {

	db, err := dbm.NewGoLevelDB(name, dbDir)
	if err != nil {
		panic(err)
	}
	state = NewState("KAppState", db).load()
	return &KApp{}
}

// 实现abci的Info协议
func (app *KApp) Info(req types.RequestInfo) (res types.ResponseInfo) {

	res.Data = cfg.Moniker
	res.LastBlockHeight = int64(state.Height)
	res.LastBlockAppHash = state.AppHash
	res.Version = req.Version

	return
}

// 实现abci的SetOption协议
func (app *KApp) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{Code: types.CodeTypeOK}
}

// 实现abci的DeliverTx协议
func (app *KApp) DeliverTx(txBytes []byte) types.ResponseDeliverTx {
	tx := NewTransaction()

	m, _ := hex.DecodeString(string(txBytes))

	// decode tx
	if err := tx.FromBytes(m); err != nil {
		return types.ResponseDeliverTx{
			Code: code.ErrTransactionDecode.Code,
			Log:  err.Error(),
		}
	}

	return types.ResponseDeliverTx{Code: code.Ok.Code}
}

// 实现abci的CheckTx协议
func (app *KApp) CheckTx(txBytes []byte) types.ResponseCheckTx {

	tx := NewTransaction()
	m, _ := hex.DecodeString(string(txBytes))

	// decode tx
	if err := tx.FromBytes(m); err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrTransactionDecode.Code,
			Log:  err.Error(),
		}
	}

	// verify sign
	if err := tx.Verify(); err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrTransactionVerify.Code,
			Log:  err.Error(),
		}
	}

	return types.ResponseCheckTx{Code: code.Ok.Code}
}

// Commit will panic if InitChain was not called
func (app *KApp) Commit() (res types.ResponseCommit) {

	appHash := make([]byte, 8)
	binary.PutVarint(appHash, state.Size)
	state.AppHash = appHash
	state.Height ++

	state.save()
	return types.ResponseCommit{Data: appHash}
}

func (app *KApp) Query(reqQuery types.RequestQuery) (res types.ResponseQuery) {
	return res
}

// Save the validators in the merkle tree
func (app *KApp) InitChain(req types.RequestInitChain) types.ResponseInitChain {

	logger.Info("InitChain")
	for _, v := range req.Validators {

		// 最高权限拥有者
		if v.Power == 10 {

			state.db.Set([]byte("__app:genesis_validator"), v.PubKey)

			app.GValidator = hex.EncodeToString(v.PubKey)
		}

		if r := app.updateValidator(v); r != nil {
			logger.Error("Error updating validators", "r", r.Error())
		}
	}
	return types.ResponseInitChain{}
}

func (app *KApp) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	app.ValUpdates = make([]types.Validator, 0)

	d := state.db.Get([]byte("__app:genesis_validator"))
	app.GValidator = hex.EncodeToString(d)

	// iavl.NewVersionedTree().Hash()
	// iavl.NewTree()

	return types.ResponseBeginBlock{}
}

func (app *KApp) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

// ---------------------------------------------

// 更新validator
func (app *KApp) updateValidator(v types.Validator) error {
	key := []byte(cnst.ValidatorPrefix + hex.EncodeToString(v.PubKey))

	// power等于-1的时候,开放节点的权限
	if v.Power == -1 {
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(&v, value); err != nil {
			return errors.New(fmt.Sprintf("Error encoding validator: %v", err))
		}

		state.db.Set(key, value.Bytes())
		state.Size ++

		logger.Info("save node ok", "key", key)

		v.Power = 0
		app.ValUpdates = append(app.ValUpdates, v)
		return nil
	}

	// power等于-2的时候,删除节点
	if v.Power == -2 {
		state.db.Delete(key)
		logger.Info("delete node ok", "key", key)

		v.Power = 0
		app.ValUpdates = append(app.ValUpdates, v)
		return nil
	}

	// power小于等于0的时候,删除验证节点
	if v.Power >= 0 {
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(&v, value); err != nil {
			return errors.New(fmt.Sprintf("Error encoding validator: %v", err))
		}

		state.db.Set(key, value.Bytes())
		state.Size ++

		logger.Info("save node ok", "key", key)

		app.ValUpdates = append(app.ValUpdates, v)
	}
	return nil
}
