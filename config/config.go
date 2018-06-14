package config

import (
	cfg "github.com/tendermint/tendermint/config"
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
		UdpPort:           8081,
		UdpHost:           "0.0.0.0",
		HttpHost:          "0.0.0.0",
		HttpPort:          8080,
		AdvertiseHttpAddr: "",
		AdvertiseUdpAddr:  "",
		LogLevel:          "debug",
	}
}
