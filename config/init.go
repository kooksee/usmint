package config

import (
	"sync"

	tlog "github.com/tendermint/tmlibs/log"
	"github.com/go-redis/redis"
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
	return a.r
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

func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		RedisUrl: "localhost:6379",
	}
}

func DefaultCfg() *Config {
	once.Do(func() {
		instance = DefaultConfig()
	})

	return instance
}
