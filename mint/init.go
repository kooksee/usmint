package mint

import (
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"
	"github.com/kooksee/usmint/config"
	"github.com/kooksee/kdb"
)

var (
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	logger log.Logger
	cfg    *config.Config
	db     *kdb.KDB
	cttm   *ContractManager
)

func Init() {
	cfg = config.DefaultCfg()
	db = cfg.App.Db()
	cttm = &ContractManager{db: db.KHash([]byte("contracts"))}
	logger = config.Log().With("module", "mint")
}
