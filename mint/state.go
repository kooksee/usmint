package mint

import (
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/types/consts"
)

func NewState(dbs ... *kdb.KDB) *State {
	db1 := db
	if len(dbs) > 0 {
		db1 = dbs[0]
	}

	name := consts.Meta(consts.StatePrefix)
	return &State{name: name, db: db1.KHash([]byte(name))}
}

type State struct {
	db   *kdb.KHash
	name string

	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func (s *State) Load() {
	stateBytes, err := s.db.Get([]byte(s.name))
	cmn.MustNotErr("state load error", err, json.Unmarshal(stateBytes, s))
}

func (s *State) Save() {
	stateBytes, err := json.Marshal(s)
	cmn.MustNotErr("state save error", err, s.db.Set([]byte(s.name), stateBytes))
}
