package web

import (
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"

	kcfg "github.com/kooksee/kchain/cfg"
	"fmt"
	"github.com/kooksee/kchain/types"
)

var (
	cfg    = kcfg.GetConfig()
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	logger log.Logger
)

func f(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

type m types.M
