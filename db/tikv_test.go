package db

import (
	"testing"
	"fmt"
	"context"
)

const pd = "101.132.96.156:2379"

func TestConnect(t *testing.T) {
	c := NewTikvStore("test", pd)
	fmt.Println(string(c.Get([]byte("hello"))))
	//for {
	//	k := []byte(uuid.New().String())
	//	c.Set(k, k)
	//	fmt.Println("c.get", string(c.Get(k)))
	//	fmt.Println("c.c", string(c.c.UUID()))
	//	c.withTxn(func(txn kv.Transaction) error {
	//		fmt.Println(txn.IsReadOnly())
	//		//fmt.Println(txn.String())
	//		fmt.Println(txn.Len())
	//		fmt.Println(txn.Size())
	//		return nil
	//	})
	//}
	txn, _ := c.c.Begin()
	defer txn.Rollback()

	p := []byte("test")
	fmt.Println(txn.IsReadOnly())
	fmt.Println(txn.Len())
	fmt.Println(txn.Size())
	//for i := 1000; i > 0; i-- {
	//	k := []byte(uuid.New().String())
	//	txn.Set(append(p, k...), []byte(""))
	//}
	fmt.Println(txn.GetMemBuffer().Size())
	fmt.Println(txn.GetMemBuffer().Len())

	iter, err := txn.Seek(p)
	if err != nil {
		panic(err.Error())
	}

	for iter.Valid() {
		fmt.Println(string(iter.Key()),string(iter.Value()))
		if err := iter.Next(); err != nil {
			panic(err.Error())
		}
	}

	if err := txn.Commit(context.TODO()); err != nil {
		panic(err.Error())
	}

	//r := util.BytesPrefix([]byte("test"))
	//iter := c.Iterator(r.Start, r.Limit)
	//for iter.Valid(){
	//	fmt.Println(string(iter.Key()))
	//	fmt.Println(string(iter.Value()))
	//	iter.Next()
	//}
}
