package usdb

import (
	"testing"
	"fmt"
	"github.com/pingcap/tidb/util"
	"github.com/pingcap/tidb/kv"
)

const pd = "101.132.96.156:2379"

func TestNamePrefix(t *testing.T) {
	c := NewTikvStore("test", pd)
	fmt.Println(string(c.Get([]byte("hello"))))

	s, _ := c.c.GetSnapshot(kv.MaxVersion)
	util.ScanMetaWithPrefix(s, []byte("test"), func(keys kv.Key, bytes []byte) bool {
		fmt.Println(string(keys))
		return keys.HasPrefix([]byte("test"))
	})
}

func TestIter(t *testing.T) {
	c := NewTikvStore("test", pd)
	fmt.Println(string(c.Get([]byte("hello"))))

	it := c.Iterator([]byte(""), []byte("123456"))
	for it.Valid() {
		fmt.Println(string(it.Key()), string(it.Value()))
		it.Next()
	}
}

func TestName1(t *testing.T) {
	Name = "test"
	Init()
	GetDb()
}

func TestConnect(t *testing.T) {
	c := NewTikvStore("test", pd)
	fmt.Println(string(c.Get([]byte("hello"))))

	//r := util.BytesPrefix([]byte("test"))
	//iter := c.Iterator(r.Start, r.Limit)
	//for iter.Valid(){
	//	fmt.Println(string(iter.Key()))
	//	fmt.Println(string(iter.Value()))
	//	iter.Next()
	//}
}
