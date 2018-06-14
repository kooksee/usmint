package reactors

//构建一个p2p网络用于数据的提前分发

import (
	"fmt"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tmlibs/log"
)

// todo: p2p网络分发

// KReactor 处理大量的TX分发
type KReactor struct {
	p2p.BaseReactor
	Name   string
	Prefix string
	ChId   byte
}

// NewKReactor
func NewKReactor() *KReactor {
	k := &KReactor{
		Name:   "KReactor",
		Prefix: "kr:",
		ChId:   byte(0x60),
	}
	k.BaseReactor = *p2p.NewBaseReactor(k.Name, k)
	return k
}

// SetLogger sets the Logger on the reactor and the underlying Mempool.
func (k *KReactor) SetLogger(l log.Logger) {
	if logger == nil {
		logger = l
	}
	k.Logger = l
}

// GetChannels implements Reactor.
// It returns the list of channels for this reactor.
func (k *KReactor) GetChannels() []*p2p.ChannelDescriptor {
	return []*p2p.ChannelDescriptor{{
		ID:                k.ChId,
		Priority:          5,
		SendQueueCapacity: 10,
	}}
}

// Receive implements Reactor.
// It adds any received transactions to the mempool.
func (k *KReactor) Receive(chID byte, src p2p.Peer, msgBytes []byte) {
	fmt.Println("dddddddd\n\n\n\n\n")
}

// 广播所有的节点
func (k *KReactor) Broadcast(msgBytes []byte) {
	k.Switch.Broadcast(k.ChId, msgBytes)
}
