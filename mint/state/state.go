package state

import (
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/kts/consts"
	"time"
	"github.com/ethereum/go-ethereum/common"
	"sync"
	"math/big"
)

type State struct {
	BlockHash []byte         `json:"block_hash"`
	Height    int64          `json:"height"`
	AppHash   []byte         `json:"app_hash"`
	Time      time.Time      `json:"time"`
	NumTxs    int32          `json:"num_txs"`
	TotalTxs  int64          `json:"total_txs"`
	Sender    common.Address `json:"sender"`
	Miner     common.Address `json:"miner"`
}

func (s *State) SetVal(hash []byte, height int) {
	cmn.MustNotErr("state SetTx error", db.Set(hash, big.NewInt(int64(height)).Bytes()))
}

func (s *State) GetTx(hash []byte) int {
	dt, err := db.Get(hash)
	if err != nil {
		cmn.MustNotErr("state GetTx error", err)
	}

	if dt == nil || len(dt) == 0 {
		return 0
	}

	return int(big.NewInt(0).SetBytes(dt).Int64())
}

// 保存状态值
func (s *State) Save() {
	stateBytes, err := cmn.JsonMarshal(s)
	cmn.MustNotErr("State Save error", err, cmn.ErrCurry(db.Set, []byte(consts.StatePrefix), stateBytes))
}

var once sync.Once
var stt *State

func GetState() *State {
	once.Do(func() {
		stt = &State{}
		stateBytes, _ := db.Get([]byte(consts.StatePrefix))
		if len(stateBytes) != 0 {
			cmn.MustNotErr("State.Load error", cmn.ErrCurry(cmn.JsonUnmarshal, stateBytes, stt))
		}
	})
	return stt
}
