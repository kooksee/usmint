package config

import (
	"sync"

	tcfg "github.com/tendermint/tendermint/config"
	tlog "github.com/tendermint/tmlibs/log"
)

var (
	once       sync.Once
	instance   *Config
	l          tlog.Logger
	configPath string
	home       string
)

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

func GetLog() tlog.Logger {
	if l == nil {
		panic("please init log")
	}
	return l
}
