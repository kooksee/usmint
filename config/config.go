package config

import (
	cfg "github.com/tendermint/tendermint/config"
	"path/filepath"
)

var (
	DefaultTendermintDir  = ".tendermint"
	defaultConfigDir      = "config"
	defaultDataDir        = "data"
	defaultConfigFileName = "config.toml"

	defaultConfigFilePath = filepath.Join(defaultConfigDir, defaultConfigFileName)
)

func DefaultCfg() *Config {
	once.Do(func() {
		instance = &Config{
			App:        DefaultAppConfig(),
			BaseConfig: cfg.DefaultBaseConfig(),
			RPC:        cfg.DefaultRPCConfig(),
			P2P:        cfg.DefaultP2PConfig(),
			Mempool:    cfg.DefaultMempoolConfig(),
			Consensus:  cfg.DefaultConsensusConfig(),
			TxIndex:    cfg.DefaultTxIndexConfig(),
		}
	})

	return instance
}

func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		RedisUrl: "localhost:6379",
	}
}
