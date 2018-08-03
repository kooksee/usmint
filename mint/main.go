package mint

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/kts"
	"github.com/kooksee/usmint/kts/code"
	"github.com/kooksee/usmint/cmn"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/abci/types"
)

func New() *Mint {
	return &Mint{
		state: NewState(),
		db:    db,
		val:   NewValidator(),
		miner: NewMiner(),
	}
}

// 创建一个bussiness层
type Mint struct {
	valUpdates []types.Validator
	state      *State
	db         kdb.IKDB
	val        *Validator
	miner      *Miner
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

// ContractDeploy 部署合约
func (m *Mint) ContractDeploy(tx *kts.Transaction) (string, error) {
	// 获取合约内容
	// 存储合约内容
	// 加载合约内容
	// 获得合约ID
	// 把合约ID返还给用户

	//tx.Data

	return "合约地址", nil
}

// ContractCall 调用合约
func (m *Mint) ContractCall() error {
	return nil
}

// ContractQuery 调用查询调用
func (m *Mint) ContractQuery() ([]byte, error) {
	return nil, nil
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
	if err := tx.VerifySign(); err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.ErrPipe("Mint CheckTx VerifySign Error", err).Error(),
		}
	}

	// 验证节点权限
	// 检查发送tx的节点有没有在区块链中,如果没有,那么该节点没有发送tx的权利
	if !m.val.Has(tx.GetPubKey()) {
		return types.ResponseCheckTx{
			Code: code.ErrInternal.Code,
			Log:  cmn.F("the node %s does not exist", tx.Pubkey),
		}
	}

	switch tx.Event {
	case "validator":
		if err := cmn.ErrPipe(
			"CheckTx validator error",
			cmn.ErrCurry(m.val.Check),
			cmn.ErrCurry(m.val.Decode, data),
			cmn.ErrCurry(m.UpdateValidators, types.Validator{PubKey: m.val.PubKey, Power: m.val.Power})); err != nil {
			return err
		}

	case "ctt.deploy":
		m.ContractDeploy(tx)
	case "ctt.call":
		m.ContractCall()
	case "ctt.query":
		m.ContractQuery()
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
	case "ctt.deploy":
	case "ctt.call.*":
	case "ctt.query.*":
	}

	// 成功之后,计算一个新的app hash
	// 根据之前的app hash计算,保证用户无法篡改数据
	m.state.AppHash = crypto.Keccak256(m.state.AppHash, data)
	return types.ResponseDeliverTx{
		Code: code.Ok.Code,
	}

	return nil
}

// BeginBlock 开始区块
func (m *Mint) BeginBlock(data types.RequestBeginBlock) error {
	// 初始化验证节点

	m.valUpdates = make([]types.Validator, 0)
	m.state.Height = data.Header.Height
	m.state.Block = data.Header.LastBlockID.Hash

	return nil
}

// EndBlock 结束区块
func (m *Mint) EndBlock(data []byte) ([]types.Validator, error) {
	return m.valUpdates, nil
}

// SetMiner 设置挖矿节点
func (m *Mint) SetMiner(v []byte, miner []byte) error {
	// 根据验证节点设置矿工
	// 每一个验证节点和非验证节点都是矿工
	// 矿工是由节点自己设置的,如果节点不想设置矿工,也没什么关系

	// 需要知道验证节点的地址,需要知道矿工的地址

	return m.miner.Set(v, miner)
}

// 查询
func (m *Mint) QueryTx(data []byte) types.ResponseQuery {
	return types.ResponseQuery{}
}
