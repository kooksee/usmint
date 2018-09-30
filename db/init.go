package db

import (
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/kooksee/usmint/cmn"
)

var appDb db.DB

func DefaultDBProvider(ctx *node.DBContext) (db.DB, error) {
	return NewRedisDB(ctx.ID, "redis url")
}

func Init() {
	var err error
	appDb, err = node.DefaultDBProvider(&node.DBContext{ID: "mint_app", Config: cmn.GetCfg()})
	cmn.MustNotErr("init app db error", err)
}
