package luas

import "github.com/kooksee/usmint/cmn"

// 合约中可以访问的内置函数

type INode interface {
	ChainID() string
	NodeID() string
	NodeMoniker() string
	BlockHeight() int64
	BlockTime() int64
	BlockAppHash() string
	BlockHash() string
	UnconfirmedNum() int64
}

type Node struct {
	INode
}

func (n *Node) ChainID() string {
	return cmn.GetNode().GenesisDoc().ChainID
}

func (n *Node) NodeID() string {
	return string(cmn.GetNode().NodeInfo().ID)
}

func (n *Node) NodeMoniker() string {
	return cmn.GetNode().NodeInfo().Moniker
}

func (n *Node) BlockHeight() int64 {
	return cmn.GetNode().BlockStore().Height()
}

func (n *Node) BlockTime() int64 {
	h := cmn.GetNode().BlockStore().Height()
	return cmn.GetNode().BlockStore().LoadBlockMeta(h).Header.Time.Unix()
}

func (n *Node) BlockAppHash() string {
	h := cmn.GetNode().BlockStore().Height()
	return cmn.GetNode().BlockStore().LoadBlockMeta(h).Header.AppHash.String()
}

func (n *Node) BlockHash() string {
	h := cmn.GetNode().BlockStore().Height()
	return cmn.GetNode().BlockStore().LoadBlockMeta(h).BlockID.Hash.String()
}

func (n *Node) UnconfirmedNum() int64 {
	return int64(cmn.GetNode().MempoolReactor().Mempool.Size())
}
