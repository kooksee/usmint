package mint

import (
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"
	"github.com/kooksee/kdb"
	"ybkchain/config"
)

var (
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	logger log.Logger
	db     *kdb.KDB
	cttm   *ContractManager
)

func Init() {
	db = cfg.App.Db()
	logger = config.Log().With("module", "mint")
}
