package mint

import (
	"github.com/tendermint/tendermint/libs/log"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"path/filepath"
)

var (
	db     kdb.IKDB
	logger log.Logger
)

func Init() {
	cfg := kdb.DefaultConfig()
	cfg.InitKdb(filepath.Join(cmn.GetCfg().DBDir(), "app_db"))
	db = cfg.GetDb()

	logger = cmn.Log().With("pkg", "mint")
}
