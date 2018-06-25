package accounts

import (
	"github.com/json-iterator/go"
	"kchain/config"
	goc "github.com/patrickmn/go-cache"
	"github.com/kooksee/kdb"
)

const Prefix = "act:"

var (
	cfg   *config.Config
	json  = jsoniter.ConfigCompatibleWithStandardLibrary
	db    *kdb.KHash
	cache *goc.Cache
)

func Init() {
	cfg = config.DefaultCfg()

	var err error
	db, err = cfg.GetDb().KHash("act")
	if err != nil {
		panic(err.Error())
	}
	cache = cfg.GetCache()
}
