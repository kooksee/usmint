package web

import (
	"github.com/json-iterator/go"
	"fmt"
	"github.com/kooksee/usmint/types"
	"github.com/kooksee/usmint/config"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	cfg  *config.Config
)

func Init() {
	cfg = config.DefaultCfg()
}

func f(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

type m types.Map
