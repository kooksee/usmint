package luas

import (
	"fmt"
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

func (f *Fields) DefInt(field string, val int) {
	fmt.Println(field, val)
}

func (f *Fields) Int(field string) int {
	return 0
}

func (f *Fields) DefFloat(field string, v float64) {
	fmt.Println(field, v)
}

func (f *Fields) Float(field string) float64 {
	return 0
}

func (f *Fields) DefString(field string, v string) {
}

func (f *Fields) String(field string) string {
	return ""
}

func (f *Fields) DefMap(field string, m map[string]interface{}) {
	fmt.Println(field, m)
}

func (f *Fields) Map(field string) lua.LValue {
	dd, _ := decodeRaw(f.l, []byte(`{"a":2,"3":{"3":3}}`))
	return dd
}

func (f *Fields) CallBack(fn func(value lua.LValue)) {
	dd, _ := decodeRaw(f.l, []byte(`{"a":2,"3":{"3":3}}`))
	fn(dd)
	fn(dd)
	fn(dd)
	fn(dd)
	fn(dd)
}

//func (f *Fields) DefDb(field string, m map[string]interface{}) {
//	fmt.Println(field, m)
//}
//
//func (f *Fields) Db(field string, m map[string]interface{}) {
//	fmt.Println(field, m)
//}
