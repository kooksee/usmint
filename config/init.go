package config

import (
	"sync"

	tcfg "github.com/tendermint/tendermint/config"
	tlog "github.com/tendermint/tmlibs/log"
	"github.com/go-redis/redis"
	"github.com/kooksee/kdb"
	"path/filepath"
)

var (
	once     sync.Once
	instance *Config
	l        tlog.Logger
)

type AppConfig struct {
	RedisUrl string `mapstructure:"redis_url"`
	r        *redis.Client
	db       *kdb.KDB
}

func (a *AppConfig) InitRedis() {
	opt, err := redis.ParseURL(a.RedisUrl)
	if err != nil {
		Log().Error("解析redis url", "err", err.Error())
		panic("")
	}

	a.r = redis.NewClient(opt)

	if err != nil {
		Log().Error("redis ping", "err", a.r.Ping().Err().Error())
		panic("")
	}
}

func (a *AppConfig) Redis() *redis.Client {
	if a.r == nil {
		a.InitRedis()
	}
	return a.r
}

type Config struct {
	// Top level options use an anonymous struct
	tcfg.BaseConfig `mapstructure:",squash"`

	// Options for services
	RPC       *tcfg.RPCConfig       `mapstructure:"rpc"`
	P2P       *tcfg.P2PConfig       `mapstructure:"p2p"`
	Mempool   *tcfg.MempoolConfig   `mapstructure:"mempool"`
	Consensus *tcfg.ConsensusConfig `mapstructure:"consensus"`
	TxIndex   *tcfg.TxIndexConfig   `mapstructure:"tx_index"`
	App       *AppConfig            `mapstructure:"app"`
}

func (a *Config) Db() *kdb.KDB {
	if a.App.db == nil {
		kdb.InitKdb(filepath.Join(a.DBDir(), "app.db"))
		a.App.db = kdb.GetKdb()
	}
	return a.App.db
}

// SetRoot sets the RootDir for all Config structs
func (cfg *Config) SetRoot(root string) *Config {
	cfg.BaseConfig.RootDir = root
	cfg.RPC.RootDir = root
	cfg.P2P.RootDir = root
	cfg.Mempool.RootDir = root
	cfg.Consensus.RootDir = root
	return cfg
}

func SetLog(l1 tlog.Logger) {
	l = l1
}

func Log() tlog.Logger {
	if l == nil {
		panic("please init log")
	}
	return l
}
