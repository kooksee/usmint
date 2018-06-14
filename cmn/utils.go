package cmn

import (
	"github.com/tendermint/tmlibs/common"
)

var Fmt = common.Fmt

func BytesTrimSpace(bs []byte) []byte {
	for i, b := range bs {
		if b != 0 {
			return bs[i : len(bs)-1]
		}
	}
	return nil
}
