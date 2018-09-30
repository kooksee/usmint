package main

import (
	"testing"
	"github.com/tendermint/tendermint/types"
	"fmt"
	"encoding/json"
	dbm "github.com/tendermint/tendermint/libs/db"
	"time"

	ttypes "github.com/tendermint/tendermint/abci/types"
	"encoding/base64"
)

func TestName11(t *testing.T) {
	db, err := dbm.NewGoLevelDB("blockstore", "d1")
	if err != nil {
		panic(err)
	}

	iter := db.DB().NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println(string(key))

		var d interface{}
		if err := json.Unmarshal(value, d); err != nil {
			panic(err.Error())
		}
		fmt.Println(d)
		fmt.Println(string(value))
		fmt.Println("\n\nok")
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}

	db.Print()
}

type BlockMeta struct {
	BlockID types.BlockID `json:"block_id"` // the block hash and partsethash
	Header  *types.Header `json:"header"`   // The block's Header
}

type Block struct {
	*types.Header                 `json:"header"`
	*types.Data                   `json:"data"`
	Evidence   types.EvidenceData `json:"evidence"`
	LastCommit *types.Commit      `json:"last_commit"`
}

func TestName221(t *testing.T) {
	fmt.Println(time.Now().UnixNano() / 100000)
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now().UnixNano())
}

type PFunc1 func(in interface{}, out *interface{}) error

type PFunc func(in interface{}) (out *interface{}, err error)

func Pipe(in interface{}, pf ... PFunc1) (out interface{}, err error) {
	for _, f := range pf {
		if err = f(in, &out); err != nil {
			return
		}
		in = out
	}
	return
}

func Pipe1(in interface{}, pf ... PFunc) (out interface{}, err error) {
	for _, f := range pf {
		out, err = f(in)
		if err != nil {
			return
		}
		in = out
	}
	return
}

func a1(in interface{}, out *interface{}) error {
	out = &in
	return nil
}

func TestName21(t *testing.T) {
	ret, err := Pipe([]string{"11", "rrr"}, a1, a1, a1, a1, a1)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(ret)
}

func TestLoadValidator(t *testing.T) {
	dt, err := base64.StdEncoding.DecodeString("")
	if err != nil {
		panic(err.Error())
	}

	val := ttypes.Ed25519Validator(dt, 10)
	fmt.Println(val.String())
	// val:
	//dt, err := base64.StdEncoding.DecodeString("")
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//pk, err := cryptoAmino.PubKeyFromBytes(dt)
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//val := types.NewValidator(pk, 10)
	//fmt.Println(val.String())
}
