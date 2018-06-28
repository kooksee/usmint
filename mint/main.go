package mint

import (
	"github.com/tendermint/abci/types"
	"github.com/kooksee/kdb"
	kts "github.com/kooksee/usmint/types"
)

func New() *Mint {
	return &Mint{
		state: NewState(),
		db:    db,
	}
}

// 创建一个bussiness层
type Mint struct {
	valUpdates []types.Validator
	state      *State
	db         *kdb.KDB
	val        *Validator
	ctt        *Contract
}

func (m *Mint) GetState() *State {
	return m.state
}

// UpdateValidators 更新Validators
func (m *Mint) UpdateValidators(Validators ... types.Validator) error {
	return m.val.UpdateValidator(nil)
}

// ContractDeploy 部署合约
func (m *Mint) ContractDeploy() error {
	return nil
}

// ContractCall 调用合约
func (m *Mint) ContractCall() error {
	return nil
}

// InitChain 初始化chain
func (m *Mint) InitChain(Validators ... types.Validator) error {
	return nil
}

// Commit 提交tx
func (m *Mint) Commit() []byte {
	return nil
}

// CheckTx 预提交
func (m *Mint) CheckTx(data []byte) error {
	tx, err := kts.DecodeTx(data)
	if err != nil {
		return err
	}

	if err := tx.VerifySign(); err != nil {
		return err
	}

	switch tx.Event {
	case "":

	}
	return nil
}

// DeliverTx 提交
func (m *Mint) DeliverTx(data []byte) error {
	tx, err := kts.DecodeTx(data)
	if err != nil {
		return err
	}

	switch tx.Event {
	case "":

	}

	return nil
}

// BeginBlock 开始区块
func (m *Mint) BeginBlock(data []byte) error {
	return nil
}

// EndBlock 结束区块
func (m *Mint) EndBlock(data []byte) ([]types.Validator, error) {
	return m.valUpdates, nil
}
