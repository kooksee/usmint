package luas

import (
	"github.com/yuin/gopher-lua"
	"github.com/go-redis/redis"
)

/*
gopher-lua redis 操作
智能合约跟redis进行交互和数据存储
 */

func NewRedis(contractAddress []byte) *RedisDb {
	return &RedisDb{}
}

type RedisDb struct {
	address []byte
	r       *redis.Client
	l       *lua.LState
}

func (db *RedisDb) Set() {
	db.r.Ping().String()
}