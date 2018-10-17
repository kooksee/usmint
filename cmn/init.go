package cmn

import (
	"github.com/kooksee/cmn"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/config"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/node"
	"golang.org/x/crypto/ripemd160"
	"fmt"
)

var ErrPipe = cmn.Err.ErrWithMsg
var ErrCurry = cmn.Err.Curry
var F = cmn.F
var Err = cmn.Err.Err

var MustNotErr = cmn.Err.MustNotErr
var JsonMarshal = cmn.Json.Marshal
var JsonMarshalToString = cmn.Json.MarshalToString
var JsonUnmarshal = cmn.Json.Unmarshal

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

func InitAppDb(path string) {
	cfg := kdb.DefaultConfig()
	cfg.DbPath = path
	cfg.InitKdb()
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
	return n
}

func InitNode(n1 *node.Node) {
	n = n1
}

func Ripemd160(bytes []byte) []byte {
	h := ripemd160.New()
	h.Write(bytes)
	return h.Sum(nil)
}

var maxMsgSize = 1024 * 1024
func CheckMsgSize(txBytes []byte) error {
	if len(txBytes) > maxMsgSize {
		return fmt.Errorf("msg size exceeds max size (%d > %d)", len(txBytes), maxMsgSize)
	}
	return nil
}


