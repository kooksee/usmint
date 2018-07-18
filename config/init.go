package config

import (
	"sync"

	tcfg "github.com/tendermint/tendermint/config"
	tlog "github.com/tendermint/tmlibs/log"
)

var (
	once     sync.Once
	instance *Config
	l        tlog.Logger
)

type AppConfig struct {
	DbPath string `mapstructure:"db_path" yaml:"db_path"`
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
