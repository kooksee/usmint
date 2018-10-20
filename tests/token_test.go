package tests

import (
	"testing"
	"time"
	"github.com/ethereum/go-ethereum/crypto"
	"compress/gzip"
	"fmt"
	"bytes"
	"github.com/vmihailenco/msgpack"
	"compress/zlib"
	"encoding/hex"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"io"

	secp256k11 "github.com/btcsuite/btcd/btcec"
	"encoding/base64"
	"github.com/kooksee/usmint/kts"
	"github.com/kooksee/usmint/mint/minter"
	"github.com/ethereum/go-ethereum/common"
)

func TestNodeJoin(t *testing.T) {
	//tx := kts.NewTransaction()
	//tx.Event = "node_manage"
	//tx.Timestamp = uint64(time.Now().Unix())
	//val := kts.Validator{
	//	Address: hex.EncodeToString(NodepriV.PubKey().Address().Bytes()),
	//	Power:   -1,
	//}
	//dt, _ := val.Encode()
	//tx.Data = hex.EncodeToString(dt)
	//println("sign_msg", hex.EncodeToString(tx.SignMsg()))
	//s, _ := NodepriV.Sign(tx.SignMsg())
	//tx.Pubkey = hex.EncodeToString(NodepriV.PubKey().Bytes())
	//tx.NodeSignature = hex.EncodeToString(s)
	//
	//txd, err := tx.Dumps()
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//println("tx", string(txd))
	//ret, err := abciClient.BroadcastTxCommit(txd)
	//jsonPrintln(ret)
}

func TestNodeLeave(t *testing.T) {
}

func TestPubK(t *testing.T) {
	//pp := types.Ed25519Validator(NodepriV.PubKey().Bytes(), 1)
	//println("pp", pp.PubKey.String())
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

func genPrivKey(rand io.Reader) secp256k1.PrivKeySecp256k1 {
	privKeyBytes := [32]byte{}
	_, err := io.ReadFull(rand, privKeyBytes[:])
	if err != nil {
		panic(err)
	}
	// crypto.CRandBytes is guaranteed to be 32 bytes long, so it can be
	// casted to PrivKeySecp256k1.
	return secp256k1.PrivKeySecp256k1(privKeyBytes)
}

func TestAddMiner(t *testing.T) {
	//node1PriV

	dd, err := hex.DecodeString(node1PriV)
	if err != nil {
		panic(err.Error())
	}

	p1, err := crypto.ToECDSA(dd)
	if err != nil {
		panic(err.Error())
	}
	sig, err := crypto.Sign(crypto.Keccak256([]byte("123")), p1)
	if err != nil {
		panic(err.Error())
	}

	ppb, err := crypto.SigToPub(crypto.Keccak256([]byte("123")), sig)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(crypto.PubkeyToAddress(*ppb).Hex())

	pp2 := genPrivKey(bytes.NewBuffer(dd))

	_, pubkeyObject := secp256k11.PrivKeyFromBytes(secp256k11.S256(), pp2[:])
	fmt.Println(hex.EncodeToString(pubkeyObject.SerializeUncompressed()))

	m2, er := base64.StdEncoding.DecodeString("E23TrD8rjCvnmXPyTpMzHU/RH+7nGz4u5T5oIFUjHq0=")
	if er != nil {
		panic(er.Error())
	}

	fmt.Println(hex.EncodeToString(m2))

	//m2k, err := cryptoAmino.PubKeyFromBytes(m2)
	//if err != nil {
	//	panic(err.Error())
	//}

	pubKeyBytes := [32]byte{}
	io.ReadFull(bytes.NewBuffer(m2), pubKeyBytes[:])
	// crypto.CRandBytes is guaranteed to be 32 bytes long, so it can be
	// casted to PrivKeySecp256k1.
	m2k := ed25519.PubKeyEd25519(pubKeyBytes)
	fmt.Println(hex.EncodeToString(m2k[:]))
}

func TestSetMiner(t *testing.T) {
	for a := 100; a > 0; a-- {
		tx := kts.NewTransaction()
		tx.Data = (&minter.SetMiner{
			Addr:  common.HexToAddress("0x2BFb20449ab700f477B3D1903D3d92DeE6518b2B"),
			Power: 10,
		}).Encode()

		fmt.Println(tx.Timestamp)

		dd, err := hex.DecodeString(node1PriV)
		if err != nil {
			panic(err.Error())
		}

		p1, err := crypto.ToECDSA(dd)
		if err != nil {
			panic(err.Error())
		}
		tx.DoNSign(p1)
		tx.DoSign(p1)
		res, err := abciClient.BroadcastTxCommit(tx.Encode())
		if err != nil {
			panic(err.Error())
		}
		jsonPrintln(res)
	}
}

func TestDeleteMiner(t *testing.T) {
	for a := 100; a > 0; a-- {
		tx := kts.NewTransaction()
		tx.Data = (&minter.DeleteMiner{
			Addr: common.HexToAddress("0x2BFb20449ab700f477B3D1903D3d92DeE6518b2B"),
		}).Encode()

		dd, err := hex.DecodeString(node1PriV)
		if err != nil {
			panic(err.Error())
		}

		p1, err := crypto.ToECDSA(dd)
		if err != nil {
			panic(err.Error())
		}
		tx.DoNSign(p1)
		tx.DoSign(p1)

		//res, err := abciClient.BroadcastTxAsync(tx.Encode())
		res, err := abciClient.BroadcastTxCommit(tx.Encode())
		if err != nil {
			panic(err.Error())
		}
		jsonPrintln(res)
	}
}
