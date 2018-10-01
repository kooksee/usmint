package reactors

import (
	"github.com/kooksee/usmint/wire"
)

func init() {
	RegisterRTx(wire.GetCodec())
}

func decodeMsg(bz []byte) (msg interface{}, err error) {
	return msg, wire.GetCodec().UnmarshalBinaryBare(bz, &msg)
}
