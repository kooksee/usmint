package luas

import (
	"github.com/json-iterator/go"
	"github.com/kooksee/kdb"
	"github.com/tendermint/tmlibs/log"
	"kchain/config"
	goc "github.com/patrickmn/go-cache"
)

const Prefix = "contract:"

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var logger log.Logger
var cfg *config.Config
var db *kdb.KHash
var cache *goc.Cache

func Init() {
	cfg = config.DefaultCfg()
	logger = cfg.GetLog().With("module", "luas")
	var err error
	db, err = cfg.GetDb().KHash("contract")
	if err != nil {
		panic(err.Error())
	}

	cache = cfg.GetCache()
}
