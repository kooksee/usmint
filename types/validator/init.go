package validator

import (
	"github.com/json-iterator/go"
	"kchain/config"
	"github.com/kooksee/kdb"
	"github.com/tendermint/tmlibs/log"
)

const Prefix = "val:"

var (
	cfg    *config.Config
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	db     *kdb.KHash
	logger log.Logger
)

func Init() {
	cfg = config.DefaultCfg()

	logger = cfg.GetLog().With("package", "validator")

	var err error
	db, err = cfg.GetDb().KHash("val")
	if err != nil {
		panic(err.Error())
	}
}
