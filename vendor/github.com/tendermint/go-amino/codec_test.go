package amino_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
)

type SimpleStruct struct {
	String string
	Bytes  []byte
	Time   time.Time
}

func newSimpleStruct() SimpleStruct {
	s := SimpleStruct{
		String: "hello",
		Bytes:  []byte("goodbye"),
		Time:   time.Now().UTC().Truncate(time.Millisecond), // strip monotonic and timezone.
	}
	return s
}

func TestMarshalUnmarshalBinaryPointer0(t *testing.T) {

	var s = newSimpleStruct()
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(s) // no indirection
	assert.Nil(t, err)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2) // no indirection
	assert.Nil(t, err)
	assert.Equal(t, s, s2)

}

func TestMarshalUnmarshalBinaryPointer1(t *testing.T) {

	var s = newSimpleStruct()
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(&s) // extra indirection
	assert.Nil(t, err)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2) // no indirection
	assert.Nil(t, err)
	assert.Equal(t, s, s2)

}

func TestMarshalUnmarshalBinaryPointer2(t *testing.T) {

	var s = newSimpleStruct()
	var ptr = &s
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(&ptr) // double extra indirection
	assert.Nil(t, err)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2) // no indirection
	assert.Nil(t, err)
	assert.Equal(t, s, s2)

}

func TestMarshalUnmarshalBinaryPointer3(t *testing.T) {

	var s = newSimpleStruct()
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(s) // no indirection
	assert.Nil(t, err)

	var s2 *SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2) // extra indirection
	assert.Nil(t, err)
	assert.Equal(t, s, *s2)
}

func TestMarshalUnmarshalBinaryPointer4(t *testing.T) {

	var s = newSimpleStruct()
	var ptr = &s
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinary(&ptr) // extra indirection
	assert.Nil(t, err)

	var s2 *SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2) // extra indirection
	assert.Nil(t, err)
	assert.Equal(t, s, *s2)

}

func TestCodecSeal(t *testing.T) {

	type Foo interface{}
	type Bar interface{}

	cdc := amino.NewCodec()
	cdc.RegisterInterface((*Foo)(nil), nil)
	cdc.Seal()

	assert.Panics(t, func() { cdc.RegisterInterface((*Bar)(nil), nil) })
	assert.Panics(t, func() { cdc.RegisterConcrete(int(0), "int", nil) })
}
