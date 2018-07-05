package mint

import (
	"github.com/tendermint/abci/types"
	"github.com/kooksee/kdb"
	kts "github.com/kooksee/usmint/types"
	"encoding/binary"
	"github.com/kooksee/usmint/cmn"
)

func New() *Mint {
	return &Mint{
		state: NewState(),
		db:    db,
		val:   NewValidator(),
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
func (m *Mint) UpdateValidators(vals ... types.Validator) error {
	for _, val := range vals {
		if err := m.val.UpdateValidator(&val); err != nil {
			return err
		}

		m.valUpdates = append(m.valUpdates, val)
	}
	return nil
}

// ContractDeploy 部署合约
func (m *Mint) ContractDeploy() error {
	return nil
}

// ContractCall 调用合约
func (m *Mint) ContractCall() error {
	return nil
}

func (m *Mint) ContractQuery() error {
	return nil
}

// InitChain 初始化chain
func (m *Mint) InitChain(vals ... types.Validator) error {
	for _, val := range vals {
		if err := m.val.UpdateValidator(&val); err != nil {
			return err
		}
	}
	return nil
}

// Commit 提交tx
func (m *Mint) Commit() []byte {
	if m.state.Size <= 0 {
		m.state.Size = 0
	}

	hash := make([]byte, 8)
	binary.BigEndian.PutUint64(hash, uint64(m.state.Size))

	m.state.Height++
	m.state.AppHash = hash

	m.state.Save()
	return m.state.AppHash
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

	pubkey := tx.GetPubKey()

	switch tx.Event {
	case "validator":
		if err := cmn.ErrPipeLog(
			"validator",
			m.val.Check(),
			m.val.Decode(data),
			m.UpdateValidators(types.Validator{PubKey: m.val.PubKey, Power: m.val.Power})); err != nil {
			return err
		}

	case "ctt.deploy":
		m.ContractDeploy()
	case "ctt.call":
		m.ContractCall()
	case "ctt.query":
		m.ContractQuery()
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
	case "node.validator":
	case "ctt.deploy":
	case "ctt.call.*":
	case "ctt.query.*":
	}

	return nil
}

// BeginBlock 开始区块
func (m *Mint) BeginBlock(data []byte) error {
	// 初始化验证节点
	m.valUpdates = make([]types.Validator, 0)
	return nil
}

// EndBlock 结束区块
func (m *Mint) EndBlock(data []byte) ([]types.Validator, error) {
	return m.valUpdates, nil
}
