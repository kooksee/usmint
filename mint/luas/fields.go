package luas

import (
	"github.com/yuin/gopher-lua"
)

const FieldsPrefix = "fields"

type Fields struct {
	address []byte
	l       *lua.LState
}

func (f *Fields) prefix() []byte {
	return append([]byte(Prefix+":"), f.address...)
}

func (f *Fields) k([]byte) []byte {
	return append([]byte(Prefix+":"), f.address...)
}

func (f *Fields) CallBack(fn func(value lua.LValue)) {
	dd, _ := decodeRaw(f.l, []byte(`{"a":2,"3":{"3":3}}`))
	fn(dd)
	fn(dd)
	fn(dd)
	fn(dd)
	fn(dd)
}

type T1 struct {
}

func (t *T1) HelloWe() string {
	return "jjj"
}

func (f *Fields) CallBack1(fn func(value *T1)) {
	t := &T1{}
	fn(t)
}

// 自定义全局类型
// 自定义数据存储
// 自定义kv存储
