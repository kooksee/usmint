package app

import dbm "github.com/tendermint/tmlibs/db"

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
	name    string `json:"name"`
}

func (s *State) load() *State {
	stateBytes := s.db.Get([]byte(s.name))
	if len(stateBytes) != 0 {
		if err := json.Unmarshal(stateBytes, s); err != nil {
			panic(err)
		}
	}
	return s
}

func (s *State) save() {
	stateBytes, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	s.db.Set([]byte(s.name), stateBytes)
}

func NewState(name string, db dbm.DB) *State {
	return &State{name: name, db: db}
}
