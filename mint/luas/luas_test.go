package luas

import (
	"testing"
	"github.com/yuin/gopher-lua"
	"fmt"
	"github.com/layeh/gopher-luar"
	"github.com/PuerkitoBio/goquery"
	"encoding/base64"
	"math/big"
	"encoding/hex"
)

func TestName(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	//f := &Fields{l: L}
	//L.SetGlobal("field", luar.New(L, f))

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
	if err := Map(L, "person", &person); err != nil {
		panic(err)
	}

	fmt.Printf("%s %d", person.Name, person.Age)
}

func TestName3(t *testing.T) {
	ff := func(d []byte) {
		fmt.Println(len(d))
	}
	ff(nil)
}

const script = `
	local b={}
	doc:Find(".picture_info .picture_body .wrapper .item .nav a"):Each(function(i , s )
		b[i]=s:Text()
	end)

	returns={
		a=doc.Url,
		b=b
	}
`

func TestName4(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	doc, err := goquery.NewDocument("http://www.meisupic.com/goods.php?media_id=89619142")
	if err != nil {
		panic(err.Error())
	}

	//doc.Find(".picture_info .picture_body .wrapper .item .nav a").Each(func(i int, selection *goquery.Selection) {
	//	fmt.Println(selection.Text())
	//})

	L.SetGlobal("doc", luar.New(L, doc))
	if err := L.DoString(script); err != nil {
		panic(err)
	}
	dt, err := LValueDumps(L.GetGlobal("returns"))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(dt))
}

func TestName11(t *testing.T) {
	//dt, err := base64.StdEncoding.DecodeString("vsvjEbxKPN6FrgnIO218XYJCt17CbPZyfUXrjbhgUoYdVZKvIpqBJi2hOHkpl31JUhiflb4POYNYPAeVMBZCuBEe2SA0uh9QHvkvEU3C4YKDiretKicH5fJcl5T27Lfr")
	//dt, err := base64.StdEncoding.DecodeString("MjM5NzY1MDMwMA==")
	// https://mp.weixin.qq.com/s/7eIsRcltbwHmQyTA33uy6w
	dt, err := base64.StdEncoding.DecodeString("7eIsRcltbwHmQyTA33uy6w==")
	//dt, err := base64.StdEncoding.DecodeString("A2CVPEMf8bKM%252FhIOAqcKboKqgdrRKOyM_fOWUw=")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(dt))
	fmt.Println(dt)
	dd := big.NewInt(0).SetBytes(dt).Int64() + 1
	fmt.Println(base64.StdEncoding.EncodeToString(big.NewInt(dd).Bytes()))
}

func TestName12(t *testing.T) {
	//doc, err := goquery.NewDocument("https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MjM5NzY1MDMwMA==&scene=124&")
	//if err != nil {
	//	panic(err.Error())
	//}
}

func TestName13(t *testing.T) {
	//	OFTEIKodlfVRORwQFu5V_A
	dt ,err:= hex.DecodeString("efee83a498aef356eeb7cc23bd0b620d")
	if err !=nil{
		panic(err.Error())
	}
	fmt.Println(string(dt))
	fmt.Println(dt)
}
