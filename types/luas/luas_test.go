package luas

import (
	"testing"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

func TestName(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	f := &Fields{l: L}
	L.SetGlobal("field", luar.New(L, f))

	testReturn(t, L, `return field:DefInt("n",1)`)
	testReturn(t, L, `return field:Int("n")`, "0")

	testReturn(t, L, `return field:DefMap("n",{a=2,["2"]=3})`)
	testReturn(t, L, `m=field:Map("n");print(m.a);return m.a`, "2")

	testReturn(t, L, `field:CallBack(function(n) print(n.a) end)`)
}
