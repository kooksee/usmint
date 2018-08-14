package kts

import (
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/ethereum/go-ethereum/common"
)

type IMint interface {
	//	节点加入,节点退出
	NodeJoin(val *types.Validator) error

	//	矿工地址设置
	MinerSet(address crypto.Address, address2 common.Address) error
	//	矿工地址删除
	MinerDel(address crypto.Address) error

	// metadata 存储
	MetaDataSet(dna []byte, data []byte) error
	MetaDataGet(dna []byte) []byte
}
