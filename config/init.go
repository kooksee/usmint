package config

import (
	"sync"

	cfg "github.com/tendermint/tendermint/config"
	tlog "github.com/tendermint/tmlibs/log"
	goc "github.com/patrickmn/go-cache"
	"time"
)

var (
	once       sync.Once
	instance   *Config
	l          tlog.Logger
	configPath string
	cache      = goc.New(time.Minute, 5*time.Minute)
	home       string
)

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

type AppConfig struct {
	Name     string
	Addr     string
	DbPath   string `mapstructure:"db_path" yaml:"db_path"`
	LogPath  string `mapstructure:"log_path" yaml:"log_path"`
	LogLevel string `mapstructure:"log_level" yaml:"log_level"`

	UdpPort int    `mapstructure:"udp_port" yaml:"udp_port"`
	UdpHost string `mapstructure:"udp_host" yaml:"udp_host"`

	HttpPort int    `mapstructure:"http_port" yaml:"http_port"`
	HttpHost string `mapstructure:"http_host" yaml:"http_host"`

	ExtIP string `mapstructure:"ext_ip" yaml:"ext_ip"`

	AdvertiseUdpAddr  string `mapstructure:"advertise_udp_addr" yaml:"advertise_udp_addr"`
	AdvertiseHttpAddr string `mapstructure:"advertise_http_addr" yaml:"advertise_http_addr"`

	Seeds []string `mapstructure:"seeds" yaml:"seeds"`
	PriV  string   `mapstructure:"priv" yaml:"priv"`
}

type Config struct {
	// Top level options use an anonymous struct
	cfg.BaseConfig `mapstructure:",squash"`

	// Options for services
	RPC       *cfg.RPCConfig       `mapstructure:"rpc"`
	P2P       *cfg.P2PConfig       `mapstructure:"p2p"`
	Mempool   *cfg.MempoolConfig   `mapstructure:"mempool"`
	Consensus *cfg.ConsensusConfig `mapstructure:"consensus"`
	TxIndex   *cfg.TxIndexConfig   `mapstructure:"tx_index"`
	App       *AppConfig           `mapstructure:"app"`
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
