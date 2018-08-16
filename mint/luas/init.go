package luas

import (
	"github.com/json-iterator/go"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/libs/log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var logger log.Logger
var db kdb.IKDB

func Init() {
	logger = cmn.Log().With("pkg", "luas")
}
