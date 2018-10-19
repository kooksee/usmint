package tests

import (
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/json-iterator/go"
	"encoding/hex"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var abciClient *client.HTTP
var clientPriV *ecdsa.PrivateKey

var node1PriV = "3721cf4e988845ddc5b5bd9e942da684990070a6a87bdc3c6fc32d1d9c2267bf"
var node2PriV = "d877a0fa7a6423f1b1e259d5e09240f3c478d96032970bfad71dd4cb880ab7b5"

var otherClientAddress = "2ee06d8ec74bf935ece278644891dfbc997735c4"

const ethPrivKey = "a168ea1e8ee6c5fb8c761369aa04d2c89a7cef40599596ef0f32909946e01a2f"

func init() {
	ddd, _ := hex.DecodeString(ethPrivKey)
	clientPriV, _ = crypto.ToECDSA(ddd)
	//add := crypto.PubkeyToAddress(clientPriV.PublicKey)

	//pvv, err := ethc.GenerateKey()
	//fmt.Println(hex.EncodeToString(ethc.FromECDSA(pvv)))

	//pvv1, err := ethc.GenerateKey()
	//fmt.Println(hex.EncodeToString(ethc.FromECDSA(pvv1)))

	//fmt.Println(hex.EncodeToString(ethc.PubkeyToAddress(pvv.PublicKey).Bytes()))

	//abciClient = client.NewHTTP("http://localhost:46657", "")
	abciClient = client.NewHTTP("http://localhost:26657", "")
}

type Map map[string]string

func (m Map) String() string {
	d, _ := json.Marshal(m)
	return string(d)
}

func DecodeBase64(d string) []byte {
	dd, _ := base64.StdEncoding.DecodeString(d)
	return dd
}

func DecodeBase64Str(d string) string {
	dd, _ := base64.StdEncoding.DecodeString(d)
	return string(dd)
}

func jsonPrintln(data interface{}) {
	fmt.Println(json.MarshalToString(data))
}

func println(msg string, data string) {
	fmt.Println(msg, data)
}
