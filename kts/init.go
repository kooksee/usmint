package kts

import "github.com/kooksee/usmint/cmn/wire"

type M map[string]interface{}

func init() {
	wire.Register("m", &M{})
}
