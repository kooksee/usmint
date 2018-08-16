package mint

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/kts"
	"github.com/kooksee/usmint/kts/code"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"encoding/hex"
)

func New() *Mint {
	return &Mint{
		state: NewState(),
		db:    db,
		val:   NewValidatorManager(),
		miner: NewMiner(),
		sc:    newContractManager(),
	}
}

// 创建一个bussiness层
type Mint struct {
	valUpdates []types.Validator
	state      *State
	db         kdb.IKDB
	val        *ValidatorManager
	miner      *Miner
	sc         *ContractManager
}

func (m *Mint) State() *State {
	return m.state
}

// InitChain 初始化chain
func (m *Mint) InitChain(vals ... types.Validator) {
	for _, val := range vals {
		cmn.MustNotErr("Mint InitChain", m.val.UpdateValidator(&kts.Validator{
			Address: hex.EncodeToString(val.Address),
			Power:   val.Power,
		}))
	}
}

// Commit 提交tx
func (m *Mint) Commit() []byte {
	return m.state.Save()
}

// CheckTx 预提交
func (m *Mint) CheckTx(data []byte) types.ResponseCheckTx {
	// 解析tx
	tx, err := kts.DecodeTx(data)
	if err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.ErrPipe("Mint CheckTx DecodeTx Error", err).Error(),
		}
	}

	// 签名验证
	if err := tx.VerifyNodeSign(); err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.ErrPipe("Mint CheckTx VerifySign Error", err).Error(),
		}
	}

	// 验证节点权限
	// 检查发送tx的节点有没有在区块链中,如果没有,那么该节点没有发送tx的权利
	if !m.val.Has(tx.GetPubkey().Address()) {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.F("the node %s does not exist", tx.Pubkey),
		}
	}

	switch tx.Event {
	case "node_manage":
		if err := m.val.CheckValidatorWithTx(tx); err != nil {
			return types.ResponseCheckTx{
				Code: code.ErrInternal.Code,
				Log:  cmn.ErrPipe("Mint CheckTx node_manage", err).Error(),
			}
		}

	case "sc_dp":
		if err := m.sc.DeployCheckWithTx(tx); err != nil {
			return types.ResponseCheckTx{
				Code: code.ErrInternal.Code,
				Log:  cmn.ErrPipe("Mint CheckTx sc_dp", err).Error(),
			}
		}

	case "sc_call":
		// 验证签名
		if err := m.sc.CallCheckWithTx(tx); err != nil {
			return types.ResponseCheckTx{
				Code: code.ErrInternal.Code,
				Log:  cmn.ErrPipe("Mint CheckTx sc_call", err).Error(),
			}
		}

	}

	return types.ResponseCheckTx{Code: code.Ok.Code}
}

// DeliverTx 提交
func (m *Mint) DeliverTx(data []byte) types.ResponseDeliverTx {
	tx, err := kts.DecodeTx(data)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.Ok.Code,
			Log:  cmn.ErrPipe("Mint DeliverTx", err).Error(),
		}
	}

	switch tx.Event {
	case "node_manage":
		if val, err := m.val.UpdateValidatorWithTx(tx); err != nil {
			return types.ResponseDeliverTx{
				Code: code.ErrInternal.Code,
				Log:  err.Error(),
			}
		} else {
			m.valUpdates = append(m.valUpdates, val)
		}

	case "sc_dp":
		if err := m.sc.DeployWithTx(tx); err != nil {
			return types.ResponseDeliverTx{
				Code: code.ErrInternal.Code,
				Log:  cmn.ErrPipe("Mint DeliverTx sc_dp", err).Error(),
			}
		}

	case "sc_call":
		if err := m.sc.CallWithOutRetWithTx(tx); err != nil {
			return types.ResponseDeliverTx{
				Code: code.ErrInternal.Code,
				Log:  cmn.ErrPipe("Mint DeliverTx sc_call", err).Error(),
			}
		}
	}

	// 成功之后,计算一个新的app hash
	// 根据之前的app hash计算,保证用户无法篡改数据
	m.state.AppHash = crypto.Ripemd160(append(m.state.AppHash, data...))
	return types.ResponseDeliverTx{Code: code.Ok.Code}
}

// BeginBlock 开始区块
func (m *Mint) BeginBlock(data types.RequestBeginBlock) error {
	// 初始化验证节点

	m.valUpdates = make([]types.Validator, 0)
	m.state.Height = data.Header.Height
	m.state.Block = data.Header.LastBlockHash

	return nil
}

// EndBlock 结束区块
func (m *Mint) EndBlock(data []byte) ([]types.Validator, error) {
	return m.valUpdates, nil
}

// 查询
func (m *Mint) QueryTx(data []byte) types.ResponseQuery {
	return types.ResponseQuery{}
}
