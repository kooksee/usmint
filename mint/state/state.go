package state

import (
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/cmn/consts"
	"time"
	"github.com/ethereum/go-ethereum/common"
	"sync"
	"github.com/kooksee/usmint/cmn/wire"
)

type State struct {
	BlockHash  []byte            `json:"block_hash"`
	Height     int64             `json:"height"`
	AppHash    []byte            `json:"app_hash"`
	Time       time.Time         `json:"time"`
	NumTxs     int32             `json:"num_txs"`
	TotalTxs   int64             `json:"total_txs"`
	Sender     common.Address    `json:"sender"`
	Miner      common.Address    `json:"miner"`
}

// 保存状态值
func (s *State) Save() {
	cmn.MustNotErr("State Save Error", db.Set([]byte(consts.StatePrefix), wire.Encode(s)))
}

var once sync.Once
var stt *State

func GetState() *State {
	once.Do(func() {
		stt = &State{}
		stateBytes, _ := db.Get([]byte(consts.StatePrefix))
		if len(stateBytes) != 0 {
			cmn.MustNotErr("State Load Error", wire.Decode(stateBytes, stt))
		}
	})
	return stt
}
