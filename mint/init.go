package mint

import (
	"github.com/tendermint/tendermint/libs/log"
	"github.com/kooksee/kdb"
)

var (
	db     *kdb.IKDB
	logger log.Logger
)

func Init(logger log.Logger) {
	kdb.DefaultConfig()
	logger = logger.With("pkg", "mint")
}
