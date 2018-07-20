package mint

import (
	"github.com/tendermint/abci/types"
	"github.com/kooksee/kdb"
	kts "github.com/kooksee/usmint/types"
	"encoding/binary"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/types/code"
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
	token      *kts.Token
}

func (m *Mint) State() *State {
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
func (m *Mint) CheckTx(data []byte) types.ResponseCheckTx {
	tx, err := kts.DecodeTx(data)
	if err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.ErrPipeLog("Mint CheckTx DecodeTx Error", err).Error(),
		}
	}

	if err := tx.VerifySign(); err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.ErrPipeLog("Mint CheckTx VerifySign Error", err).Error(),
		}
	}

	//pubkey := tx.GetPubKey()

	// 验证签名
	// 加载状态

	switch tx.Event {
	case "validator":
		if err := cmn.ErrPipeLog(
			"Mint CheckTx validator error",
			m.val.Check(),
			m.val.Decode(data)); err != nil {
			return types.ResponseCheckTx{
				Code: code.ErrInternal.Code,
				Log:  err.Error(),
			}
		}

		// 纯粹的存储，没有任何的逻辑
	case "store":

	case "tk.TransferTo":
	case "tk.TransferFrom":
	case "tk.Approve":

	}

	return nil
}

// DeliverTx 提交
func (m *Mint) DeliverTx(data []byte) types.ResponseDeliverTx {
	tx, err := kts.DecodeTx(data)
	if err != nil {
		return err
	}

	switch tx.Event {
	case "node.validator":

	}

	return nil
}

// 查询
func (m *Mint) QueryTx(data []byte) types.ResponseQuery {
	tx, err := kts.DecodeTx(data)
	if err != nil {
		return err
	}

	switch tx.Event {
	case "tk.TotalSupply":
	case "tk.Balances":
	case "tk.BalanceOf":
	case "tk.Allowance":
	}

	return nil
}

// BeginBlock 开始区块
func (m *Mint) BeginBlock(data types.RequestBeginBlock) error {
	// 初始化验证节点
	m.valUpdates = make([]types.Validator, 0)
	return nil
}

// EndBlock 结束区块
func (m *Mint) EndBlock(data types.RequestEndBlock) ([]types.Validator, error) {
	return m.valUpdates, nil
}

// UpdateValidator 更新验证节点,添加或者删除挖矿节点
func (m *Mint) UpdateValidator(val *types.Validator) error {
	// 其他节点的接入需要有主帐号控制
	// 主帐号只控制节点的接入退出,但是并不能控制节点的币的操作
	return m.val.UpdateValidator(val)
}
