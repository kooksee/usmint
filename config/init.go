package config

import (
	"sync"

	tcfg "github.com/tendermint/tendermint/config"
	tlog "github.com/tendermint/tmlibs/log"
	"github.com/go-redis/redis"
	"github.com/kooksee/usmint/cmn"
)

var (
	once     sync.Once
	instance *Config
	l        tlog.Logger
)

type AppConfig struct {
	RedisUrl string `mapstructure:"redis_url"`
	r        *redis.Client
}

func (a *AppConfig) InitRedis() {
	opt, err := redis.ParseURL(a.RedisUrl)
	cmn.MustNotErr("解析redis url", err)

	a.r = redis.NewClient(opt)

	cmn.MustNotErr("redis ping", a.r.Ping().Err())
}

func (a *AppConfig) Redis() *redis.Client {
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
