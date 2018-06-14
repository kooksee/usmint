package cmn

import (
	"github.com/json-iterator/go"
	"github.com/kooksee/usmint/config"
	"github.com/tendermint/tmlibs/log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var logger log.Logger
var cfg *config.Config

func Init() {
	cfg = config.DefaultCfg()
	logger = config.GetLog().With("module", "cmn")
}
