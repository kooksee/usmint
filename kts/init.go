package kts

import "github.com/kooksee/usmint/cmn"

type M map[string]interface{}

func (m M) String() string {
	d, _ := cmn.JsonMarshal(m)
	return string(d)
}
