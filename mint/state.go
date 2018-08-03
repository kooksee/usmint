package mint

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/kts/consts"
	"encoding/json"
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
	cmn.MustNotErr(cmn.ErrPipe("state load error", err, cmn.Wrap(json.Unmarshal, stateBytes, s)))
}

// 保存状态值
func (s *State) Save() []byte {
	stateBytes, err := json.Marshal(s)
	cmn.MustNotErr(cmn.ErrPipe("state save error", err, s.db.Set([]byte(s.name), stateBytes)))
	return stateBytes
}
