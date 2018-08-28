package tests

import (
	"testing"
	"github.com/kooksee/usmint/kts"
	"time"
	"encoding/hex"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/ethereum/go-ethereum/crypto"
	"compress/gzip"
	"fmt"
	"bytes"
	"github.com/vmihailenco/msgpack"
	"compress/zlib"
)

func TestNodeJoin(t *testing.T) {
	tx := kts.NewTransaction()
	tx.Event = "node_manage"
	tx.Timestamp = uint64(time.Now().Unix())
	val := kts.Validator{
		Address: hex.EncodeToString(NodepriV.PubKey().Address().Bytes()),
		Power:   -1,
	}
	dt, _ := val.Encode()
	tx.Data = hex.EncodeToString(dt)
	println("sign_msg", hex.EncodeToString(tx.SignMsg()))
	s, _ := NodepriV.Sign(tx.SignMsg())
	tx.Pubkey = hex.EncodeToString(NodepriV.PubKey().Bytes())
	tx.NodeSignature = hex.EncodeToString(s)

	txd, err := tx.Dumps()
	if err != nil {
		panic(err.Error())
	}

	println("tx", string(txd))
	ret, err := abciClient.BroadcastTxCommit(txd)
	jsonPrintln(ret)
}

func TestNodeLeave(t *testing.T) {
}

func TestPubK(t *testing.T) {
	pp := types.Ed25519Validator(NodepriV.PubKey().Bytes(), 1)
	println("pp", pp.PubKey.String())
}

func TestName33(t *testing.T) {
	data := []byte("hello")

	var ds [][]byte
	for i := 0; i < 10000; i++ {
		data = crypto.Keccak256(data)
		ds = append(ds, data)
	}

	tt := time.Now().UnixNano()
	ddd, _ := json.Marshal(ds)
	fmt.Println("gzip size raw:", len(ddd), time.Now().UnixNano()-tt)

	dd3, _ := msgpack.Marshal(ds)
	fmt.Println("gzip size,msgpack:", len(dd3), time.Now().UnixNano()-tt)

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(ddd)
	defer w.Close()
	w.Flush()
	fmt.Println("gzip size gzip:", len(b.Bytes()), time.Now().UnixNano()-tt)

	var b1 bytes.Buffer
	w1 := zlib.NewWriter(&b1)
	defer w1.Close()
	w1.Write(ddd)
	w1.Flush()

	fmt.Println("gzip size zlib:", len(b1.Bytes()), time.Now().UnixNano()-tt)

	//r, _ := gzip.NewReader(&b)
	//defer r.Close()
	//undatas, _ := ioutil.ReadAll(r)
	//fmt.Println("ungzip size:", len(undatas))
}
