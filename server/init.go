package server

import "github.com/tendermint/tmlibs/log"

var (
	logger log.Logger
)

func SetLogger(l log.Logger) {
	logger = l
}
