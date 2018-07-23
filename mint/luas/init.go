package luas

import (
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"
	"github.com/kooksee/usmint/config"
	"github.com/kooksee/kdb"
)

const Prefix = "contract:"

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var logger log.Logger
var cfg *config.Config
var db *kdb.KDB

func Init() {
	cfg = config.DefaultCfg()
	db = cfg.App.Db()
	logger = config.Log().With("module", "luas")
}
