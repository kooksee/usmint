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

func SetLogger(l log.Logger) {
	logger = l.With("module", "kapp")
}

func SetCfg(c *config.Config) {
	cfg = c
}
