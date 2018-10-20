package minter

import (
	"github.com/tendermint/tendermint/abci/types"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/kts"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
	"math/big"
	"github.com/kooksee/usmint/cmn/wire"
)

type Miner struct {
	Addr  common.Address
	Power int64
}

func ExistMiner(addr common.Address) bool {
	dt, err := db.Get(addr.Bytes())
	if err != nil {
		cmn.MustNotErr("ExistMiner Error", err)
	}
	return len(dt) != 0
}

func IsMaster(addr common.Address) bool {
	dt, err := db.Get(addr.Bytes())
	if err != nil {
		cmn.MustNotErr("IsMaster Error", err)
	}
	return big.NewInt(0).SetBytes(dt).Int64() == 10
}

func InitMaster() {
	cmn.Log().Error("InitMaster", "address", "0x2BFb20449ab700f477B3D1903D3d92DeE6518b2B")
	cmn.MustNotErr("InitMaster",
		db.Set(common.HexToAddress("0x2BFb20449ab700f477B3D1903D3d92DeE6518b2B").Bytes(), big.NewInt(10).Bytes()))
}

// 设置矿工地址
// v 验证节点
// ma 矿工地址
type SetMiner struct {
	*kts.BaseDataHandler
	Addr  common.Address
	Power int64
}

func (t *SetMiner) Encode() []byte {
	return wire.Encode(t)
}

func (t *SetMiner) Decode(dt []byte) error {
	return wire.Decode(dt, t)
}

func (t *SetMiner) OnCheck(tx *kts.Transaction, res *types.ResponseCheckTx) {
	if !IsMaster(tx.GetMiner()) {
		res.Code = 1
		res.Log = "you do not have permission to operate the api"
		return
	}

	if t.Power > 10 {
		res.Code = 1
		res.Log = fmt.Sprintf("the miner power must less than 10, now(%d)", t.Power)
		return
	}
}

func (t *SetMiner) OnDeliver(tx *kts.Transaction, res *types.ResponseDeliverTx) {
	if err := db.Set(t.Addr.Bytes(), big.NewInt(t.Power).Bytes()); err != nil {
		res.Code = 1
		res.Log = err.Error()
	}
}

// 删除矿工地址
// v 验证节点
type DeleteMiner struct {
	*kts.BaseDataHandler
	Addr common.Address
}

func (t *DeleteMiner) OnCheck(tx *kts.Transaction, res *types.ResponseCheckTx) {
	if !IsMaster(tx.GetMiner()) {
		res.Code = 1
		res.Log = "you do not have permission to operate the api"
		return
	}
}

func (t *DeleteMiner) OnDeliver(tx *kts.Transaction, res *types.ResponseDeliverTx) {
	if err := cmn.ErrPipe("Miner.Delete", db.Del(t.Addr.Bytes())); err != nil {
		res.Code = 1
		res.Log = err.Error()
	}
}

func (t *DeleteMiner) Encode() []byte {
	return wire.Encode(t)
}

func (t *DeleteMiner) Decode(dt []byte) error {
	return wire.Decode(dt, t)
}
