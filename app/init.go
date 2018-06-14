package app

import (
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"
	"github.com/kooksee/usmint/config"
)

var (
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	logger log.Logger
	state  *State
	cfg    *config.Config
)

func Init() {
	cfg = config.DefaultCfg()
	logger = config.GetLog().With("module", "kapp")
}
