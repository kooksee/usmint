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
		miner: NewMiner(),
	}
}

// 创建一个bussiness层
type Mint struct {
	valUpdates []types.Validator
	state      *State
	db         *kdb.KDB
	val        *Validator
	miner      *Miner
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

	//pubkey := tx.GetPubKey()

	// 验证签名
	// 检测合约是否在缓存当中,没有的就加载进来
	// 加载lua类库
	// 加载状态

	switch tx.Event {
	case "validator":
		if err := cmn.ErrPipeLog(
			"CheckTx validator error",
			m.val.Check(),
			m.val.Decode(data),
			m.UpdateValidators(types.Validator{PubKey: m.val.PubKey, Power: m.val.Power})); err != nil {
			return err
		}

		//	设置矿工的地址
	case "miner":

		//	投票智能合约
	case "vote":

	case "db.mSet":
	case "db.mSet":
	case "db.mSet":
	case "db.mSet":
	case "db.mSet":

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

// SetMiner 设置挖矿节点
func (m *Mint) SetMiner(v []byte, miner []byte) error {
	// 根据验证节点设置矿工
	// 每一个验证节点和非验证节点都是矿工
	// 矿工是由节点自己设置的,如果节点不想设置矿工,也没什么关系

	// 需要知道验证节点的地址,需要知道矿工的地址

	return m.miner.Set(v, miner)
}
