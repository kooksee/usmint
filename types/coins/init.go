package coin

import (
	"kchain/config"
	"github.com/kooksee/kdb"
)

const Prefix = "coin:"

var (
	cfg *config.Config
	db  *kdb.KHash
)

func Init() {
	cfg = config.DefaultCfg()

	var err error
	db, err = cfg.GetDb().KHash("coin")
	if err != nil {
		panic(err.Error())
	}
}
