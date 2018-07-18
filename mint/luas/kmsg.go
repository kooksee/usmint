package luas

func NewMsg() *Msg {
	return &Msg{}
}

type Msg struct {
}

// 调用数据
func (m *Msg) Data() []byte {
	return nil
}

// 发送地址
func (m *Msg) Sender() []byte {
	return nil
}
