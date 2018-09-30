package reactors

//构建一个p2p网络用于数据的提前分发

import (
	"fmt"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/kooksee/kdb"
	"github.com/tendermint/go-amino"
	"reflect"
)

type RTx interface{}

func RegisterRTx(cdc *amino.Codec) {
	cdc.RegisterInterface((*RTx)(nil), nil)
	cdc.RegisterConcrete(&sentTxReq{}, "tendermint/rtx/sentTxReq", nil)
	cdc.RegisterConcrete(&getTxReq{}, "tendermint/rtx/getTxReq", nil)
	cdc.RegisterConcrete(&getTxResp{}, "tendermint/rtx/getTxResp", nil)
}

type sentTxReq struct {
	Hash []byte
	Data []byte
}

type getTxReq struct {
	Hash []byte
}

type getTxResp struct {
	Hash []byte
	Data []byte
}

func decodeMsg(bz []byte) (msg RTx, err error) {
	err = cdc.UnmarshalBinaryBare(bz, &msg)
	return
}

// todo: p2p网络分发

// KReactor 处理大量的TX分发
type KReactor struct {
	p2p.BaseReactor
	Name   string
	Prefix string
	ChId   byte
	db     kdb.IKHash
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

func (k *KReactor) Receive(chID byte, src p2p.Peer, msgBytes []byte) {
	msg, err := decodeMsg(msgBytes)
	if err != nil {
		k.Logger.Error("Error decoding message", "src", src, "chId", chID, "msg", msg, "err", err, "bytes", msgBytes)
		k.Switch.StopPeerForError(src, err)
		return
	}

	k.Logger.Debug("Receive", "src", src, "chID", chID, "msg", msg)

	switch msg := msg.(type) {

	case *sentTxReq:
		b, _ := k.db.Exist(msg.Hash)
		if b {
			return
		}

		k.db.Set(msg.Hash, msg.Data)
		k.Switch.Broadcast(k.ChId, msgBytes)

	case *getTxReq:
		dt, _ := k.db.Get(msg.Hash)
		src.TrySend(chID, dt)
	case *getTxResp:
		k.db.Set(msg.Hash, msg.Data)
	default:
		k.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
	}

	fmt.Println("dddddddd\n\n\n\n\n")
	k.Logger.Error("dddddddd\n\n\n\n\n")
	k.Logger.Error(string(msgBytes))

}

// 广播所有的节点
func (k *KReactor) Broadcast(msgBytes []byte) {
	k.Switch.Broadcast(k.ChId, msgBytes)
}

func (k *KReactor) BroadcastTx(req sentTxReq) {
	k.Switch.Broadcast(k.ChId, amino.MustMarshalBinaryBare(req))
}

// 广播所有的节点, 获取所有的tx数据
func (k *KReactor) GetTxReq(req getTxReq) {
	k.Switch.Broadcast(k.ChId, amino.MustMarshalBinaryBare(req))
}

// 检查tx是否存在
func (k *KReactor) Exist(txid []byte) bool {
	b, _ := k.db.Exist(txid)
	return b
}
