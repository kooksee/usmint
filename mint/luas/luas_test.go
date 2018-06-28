package luas

import (
	"testing"
	"github.com/yuin/gopher-lua"
	"github.com/layeh/gopher-luar"
	"fmt"
)

func TestName(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	f := &Fields{l: L}
	L.SetGlobal("field", luar.New(L, f))

	//testReturn(t, L, `return field:DefInt("n",1)`)
	//testReturn(t, L, `return field:Int("n")`, "0")
	//
	//testReturn(t, L, `return field:DefMap("n",{a=2,["2"]=3})`)
	//testReturn(t, L, `m=field:Map("n");print(m.a);return m.a`, "2")

	testReturn(t, L, `field:CallBack(function(n) print(n.a) end)`)
	testReturn(t, L, `field:CallBack1(function(n) print(n:hello_we()) end)`)
}

func TestInit1(t *testing.T) {
	l := lua.NewState()
	defer l.Close()
	if err := l.DoFile("main.lua"); err != nil {
		panic(err.Error())
	}

	fmt.Println(l.GetGlobal("m").String())
	fmt.Println(l.GetGlobal("m1").String())

	if t, ok := l.GetGlobal("init").(*lua.LTable); ok {
		fmt.Println(json.MarshalToString(t))
	}
	fmt.Println(l.Type())

	//testReturn(t, l, `field:CallBack(function(n) print(n.a) end)`)
}

func TestName2(t *testing.T) {
	type Role struct {
		Name string
	}

	type Person struct {
		Name      string
		Age       int
		WorkPlace string
		Role      []*Role
	}

	L := lua.NewState()
	if err := L.DoString(`
person = {
  name = "Michel",
  age  = "31", -- weakly input
  work_place = "San Jose",
  role = {
    {
      name = "Administrator"
    },
    {
      name = "Operator"
    }
  }
}
`); err != nil {
		panic(err)
	}
	var person Person
	if err := Map(L,"person", &person); err != nil {
		panic(err)
	}
	fmt.Printf("%s %d", person.Name, person.Age)
}

func TestName3(t *testing.T) {
	ff:= func(d []byte) {
		fmt.Println(len(d))
	}
	ff(nil)

}