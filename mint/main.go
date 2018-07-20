package mint

import (
	"github.com/tendermint/abci/types"
	"github.com/kooksee/kdb"
	kts "github.com/kooksee/usmint/types"
	"encoding/binary"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/types/code"
	"github.com/ethereum/go-ethereum/common"
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

		//	设置矿工的地址
	case "miner.set":

		//	删除矿工
	case "miner.del":

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

// SetMiner 设置挖矿节点
func (m *Mint) SetMiner(v common.Address, miner []byte) error {
	// 根据验证节点设置矿工
	// 每一个验证节点和非验证节点都是矿工
	// 矿工是由节点自己设置的,如果节点不想设置矿工,也没什么关系

	// 需要知道验证节点的地址,需要知道矿工的地址
	// 当然，也可以考虑放入其他的信息，以便对节点进行监督

	return m.miner.Set(v.Bytes(), miner)
}

// DelMiner 删除挖矿节点
func (m *Mint) DelMiner(v common.Address) error {
	// 删除时候，那么该账号管理的数据或者代币会回归到系统当中
	return m.miner.Delete(v.Bytes())
}

// 激励分发,用来给挖矿的节点分发激励,这个出发可以是多条件的
func (m *Mint) Inflate() error {
	// 获得所有的挖矿节点
	// 获得挖矿的贡献
	// 计算贡献
	// 然后构建一个tx,并发送tx到所有的节点
	// 节点收到了tx之后,然后保存激励奖金,然后清空自己的贡献数值
	// 关于检查自己的锁定到期，以及锁定解锁，需要节点自己去出发,如果节点不触发，那么奖励的钱不回参与币龄的计算
	// 触发可以是任何的节点,但是什么时候发放需要有主节点来控制
	return nil
}

// 奖励扣除
func (m *Mint) RewardDel() error {
	// 在节点发生作弊的时候，那么需要扣除节点的抵押代币
	return nil
}

// UpdateValidator 更新验证节点,添加或者删除挖矿节点
func (m *Mint) UpdateValidator(val *types.Validator) error {
	// 其他节点的接入需要有主帐号控制
	// 主帐号只控制节点的接入退出,但是并不能控制节点的币的操作
	return m.val.UpdateValidator(val)
}
