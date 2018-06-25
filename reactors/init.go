package reactors

import (
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"

	"github.com/kooksee/usmint/config"
)

var (
	cfg    *config.Config
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	logger log.Logger
)

func Init() {

}
func SetLogger(l log.Logger) {
	logger = l
}

func SetCfg(c *config.Config) {
	cfg = c
}
