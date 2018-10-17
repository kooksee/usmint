package wire

import (
	"github.com/tendermint/go-amino"
)

var cdc *amino.Codec

func init() {
	cdc = amino.NewCodec()
}

func Register(name string, o interface{}) {
	cdc.RegisterConcrete(o, "mint/"+name, nil)
}

func RegisterInterface(o interface{}) {
	cdc.RegisterInterface(o, nil)
}

func Encode(o interface{}) []byte {
	dt, err := cdc.MarshalBinaryBare(o)
	if err != nil {
		panic(err.Error())
	}
	return dt
}

func Decode(bz []byte, ptr interface{}) error {
	return cdc.UnmarshalBinaryBare(bz, ptr)
}
