package cmn

import (
	"github.com/kooksee/cmn"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
)

var ErrPipe = cmn.Err.ErrWithMsg
var ErrCurry = cmn.Err.Curry
var F = cmn.F
var Err = cmn.Err.Err
var StrJoin = cmn.StrJoin

var MustNotErr = cmn.Err.MustNotErr
var JsonMarshal = cmn.Json.Marshal
var JsonUnmarshal = cmn.Json.Unmarshal

type M map[string]interface{}

func (m M) String() string {
	d, _ := JsonMarshal(m)
	return string(d)
}

var logger log.Logger

func InitLog(logger1 log.Logger) {
	logger = logger1
}

func Log() log.Logger {
	if logger == nil {
		panic("please init logger")
	}
	return logger
}

var cfg *config.Config

func InitCfg(cfg1 *config.Config) {
	cfg = cfg1
}

func GetCfg() *config.Config {
	if cfg == nil {
		panic("please init config")
	}
	return cfg
}

var n *node.Node

func GetNode() *node.Node {
	if n == nil {
		panic("please init node")
	}
	n.MempoolReactor()
	return n
}

func InitNode(n1 *node.Node) {
	n = n1
}
