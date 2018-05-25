package config

import (
	cfg "github.com/tendermint/tendermint/config"
)

func GetCfg() *Config {
	if instance == nil {
		panic("please init config")
	}
	return instance
}

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
