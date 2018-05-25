package reactors

import (
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tmlibs/log"
	"fmt"
)

// KReactor 处理大量的TX分发
type KReactor struct {
	p2p.BaseReactor

	Name     string
	DbPrefix string
	ChId     byte

	peers []p2p.Peer
}

// NewKReactor
func NewKReactor() *KReactor {
	k := &KReactor{
		Name:     "KReactor",
		DbPrefix: "kr:",
		ChId:     byte(0x60),
	}
	k.BaseReactor = *p2p.NewBaseReactor(k.Name, k)
	return k
}

// SetLogger sets the Logger on the reactor and the underlying Mempool.
func (k *KReactor) SetLogger(l log.Logger) {
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

func (k *KReactor) AddPeer(peer p2p.Peer) {
	k.peers = append(k.peers, peer)
	//logger.Error(peer.String())
	//logger.Error(peer.NodeInfo().String())
}

// Receive implements Reactor.
// It adds any received transactions to the mempool.
func (k *KReactor) Receive(chID byte, src p2p.Peer, msgBytes []byte) {
	fmt.Println("dddddddd\n\n\n\n\n")
	logger.Error("pppl")
	if chID != k.ChId {
		return
	}

	logger.Error(src.String())
	logger.Error(src.NodeInfo().String())
	logger.Error(string(msgBytes))
	return

	//kda := map[string]interface{}{}
	//err := json.Unmarshal(msgBytes, &kda)
	//if err != nil {
	//	logger.Error(string(msgBytes))
	//	logger.Error(err.Error())
	//	return
	//}
	//
	//tx := &Transaction{}
	//json.Unmarshal(msgBytes, tx)
	//if tx.IsCTTOpt() {
	//	ss := NewContractManager(3333)
	//	ss.AddContract(tx.Address)
	//}
}

// 广播所有的节点
func (k *KReactor) Broadcast(msgBytes []byte) {
	for _, p := range k.peers {
		fmt.Println(p.Send(k.ChId, msgBytes))
	}
	//k.Switch.Broadcast(k.ChId, msgBytes)
}

// 二度Broadcast,消息传播二次就停止
func (k *KReactor) Broadcast2(msgBytes []byte) {
	k.Switch.Broadcast(k.ChId, msgBytes)
}
