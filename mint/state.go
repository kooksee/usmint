package mint

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/kts/consts"
)

func NewState(dbs ... kdb.IKDB) *State {
	db1 := db
	if len(dbs) > 0 {
		db1 = dbs[0]
	}

	name := consts.Meta(consts.StatePrefix)
	return &State{name: name, db: db1.KHash([]byte(name))}
}

type State struct {
	db   kdb.IKHash
	name string

	Block   []byte `json:"block"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func (s *State) Load() {
	stateBytes, err := s.db.Get([]byte(s.name))
	cmn.MustNotErr("State.Load error", err, cmn.ErrCurry(cmn.JsonUnmarshal, stateBytes, s))
}

// 保存状态值
func (s *State) Save() []byte {
	stateBytes, err := cmn.JsonMarshal(s)
	cmn.MustNotErr("State.Save error", err, cmn.ErrCurry(s.db.Set, []byte(s.name), stateBytes))
	return stateBytes
}
