package luas

func NewMsg() *Msg {
	return &Msg{}
}

type Msg struct {
}

// 完整的calldata
func (m *Msg) Data() []byte {
	return nil
}

// ﻿消息的发送者（当前调用）
func (m *Msg) Sender() []byte {
	return nil
}
