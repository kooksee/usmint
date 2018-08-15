package mint

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/kts/consts"
)

func NewState() *State {
	name := consts.Meta(consts.StatePrefix)
	return &State{db: db.KHash([]byte(name)),Height:100}
}

type State struct {
	db      kdb.IKHash
	Block   []byte `json:"block"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func (s *State) Load() {
	stateBytes, err := s.db.Get([]byte(consts.StatePrefix))
	cmn.MustNotErr("State.Load error", err, cmn.ErrCurry(cmn.JsonUnmarshal, stateBytes, s))
}

// 保存状态值
func (s *State) Save() []byte {
	stateBytes, err := cmn.JsonMarshal(s)
	cmn.MustNotErr("State.Save error", err, cmn.ErrCurry(s.db.Set, []byte(consts.StatePrefix), stateBytes))
	return stateBytes
}
