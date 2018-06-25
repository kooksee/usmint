package state

import (
	"github.com/json-iterator/go"
	"kchain/config"
	"kchain/cmn"
	"github.com/kooksee/kdb"
)

var (
	cfg  *config.Config
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	errs = cmn.Errs
	db   *kdb.KHash
)

func Init() {
	cfg = config.DefaultCfg()

	var err error
	db, err = cfg.GetDb().KHash("state")
	if err != nil {
		panic(err.Error())
	}
}
