package state

import "github.com/tendermint/go-crypto"

const Prefix = "state"

type State struct {
	Size             int64  `json:"size"`
	Height           int64  `json:"height"`
	AppHash          []byte `json:"app_hash"`
	genesisValidator crypto.PubKey
}

func (s *State) SetMasterValidator(genesisValidator []byte) error {
	pk, err := crypto.PubKeyFromBytes(genesisValidator)
	if err != nil {
		return err
	}
	s.genesisValidator = pk
	return nil
}

func (s *State) GetMasterValidator(genesisValidator []byte) crypto.PubKey {
	return s.genesisValidator
}

func LoadState() *State {
	stateBytes, err := db.Get([]byte(Prefix))
	if err != nil {
		panic(err.Error())
	}

	state := &State{}
	if len(stateBytes) != 0 {
		if err := json.Unmarshal(stateBytes, state); err != nil {
			panic(errs("loadState error", err.Error()))
		}
	}

	return state
}

func SaveState(state *State) error {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(errs("saveState error", err.Error()))
	}
	return db.Set([]byte(Prefix), stateBytes)
}
