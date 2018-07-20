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
		m.state.Size = 1
	}

	hash := make([]byte, 8)
	binary.BigEndian.PutUint64(hash, uint64(m.state.Size))

	m.state.Size++
	m.state.AppHash = hash

	m.state.Save()
	return m.state.AppHash
}

// CheckTx 预提交
func (m *Mint) CheckTx(data []byte) types.ResponseCheckTx {
	// 解析tx
	tx, err := kts.DecodeTx(data)
	if err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.ErrPipeLog("Mint CheckTx DecodeTx Error", err).Error(),
		}
	}

	// 签名验证
	if err := tx.VerifySign(); err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.ErrPipeLog("Mint CheckTx VerifySign Error", err).Error(),
		}
	}

	// 验证节点权限
	// 检查发送tx的节点有没有在区块链中,如果没有,那么该节点没有发送tx的权利
	pk, _ := tx.GetPubKey()
	if !m.val.Has(pk) {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.Fmt("the node %s does not exist", tx.Pubkey),
		}
	}

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

	}

	return types.ResponseCheckTx{
		Code: code.Ok.Code,
	}
}

// DeliverTx 提交
func (m *Mint) DeliverTx(data []byte) types.ResponseDeliverTx {
	tx, _ := kts.DecodeTx(data)

	switch tx.Event {
	case "validator":
		val := &Validator{}
		val.Decode(tx.Data)
		if err := m.UpdateValidators(types.Validator{PubKey: val.PubKey, Power: val.Power}); err != nil {
			return types.ResponseDeliverTx{
				Code: code.ErrInternal.Code,
				Log:  cmn.ErrPipeLog("DeliverTx Validator UpdateValidators Error", err).Error(),
			}
		}

	case "store":
	}

	return types.ResponseDeliverTx{
		Code: code.Ok.Code,
	}
}

// 查询
func (m *Mint) QueryTx(data []byte) types.ResponseQuery {
	return types.ResponseQuery{}
}

// BeginBlock 开始区块
func (m *Mint) BeginBlock(data types.RequestBeginBlock) error {
	// 初始化验证节点
	m.valUpdates = make([]types.Validator, 0)
	m.state.Height = data.Header.Height
	return nil
}

// EndBlock 结束区块
func (m *Mint) EndBlock(data types.RequestEndBlock) ([]types.Validator, error) {
	return m.valUpdates, nil
}
