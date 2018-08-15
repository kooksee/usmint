package kts

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/ethereum/go-ethereum/common"
)

type IMint interface {
	// 节点加入,节点退出
	NodeManage(data string) error

	// 节点的查询
	NodeQuery(addr []byte) ([]byte, error)

	// 得到所有的节点
	NodeAll() ([][]byte, error)

	// 矿工地址设置,绑定eth地址到tendermint地址
	MinerSet(address crypto.Address, address2 common.Address) error

	// 获得矿工绑定的tendermint地址
	MinerGet(address crypto.Address) (common.Address, error)

	// 矿工地址删除
	MinerDel(address crypto.Address) error

	// 合约部署
	ContractDeploy(data []byte) error

	// 合约调用
	ContractCall(cAddr []byte, method string, args ... []byte) ([]byte, error)

	// metadata 存储
	MetaDataSet(dna []byte, data []byte) error
	MetaDataGet(dna []byte) []byte
}
