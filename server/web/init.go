package web

import (
	"github.com/json-iterator/go"
	"fmt"
	"github.com/kooksee/usmint/types"
	"github.com/tendermint/tmlibs/log"
	"github.com/kooksee/usmint/config"
)

var (
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	logger log.Logger
	cfg    *config.Config
)

func Init() {
	cfg = config.DefaultCfg()
	logger = config.Log().With("module", "web")
}

func f(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

type m types.Map
