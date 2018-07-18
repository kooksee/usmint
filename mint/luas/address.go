package luas

import "encoding/hex"

func Addr(addr string) *Address {
	addr1, _ := hex.DecodeString(addr)
	return &Address{addr: addr1}
}

type Address struct {
	addr []byte
}

// 地址余额。
func (a *Address) Balance() uint64 {
	return 0
}

// ﻿向 地址类型 发送数量为 amount 的 Wei，失败时抛出异常，发送 2300 gas 的矿工费，不可调节。
func (a *Address) Transfer(addr []byte, amount uint64) error {
	return nil
}

func (a *Address) Call(mth string, args interface{}) error {
	return nil
}
