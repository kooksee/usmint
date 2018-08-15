package tests

import (
	"testing"
	"github.com/kooksee/usmint/kts"
	"time"
	"encoding/hex"
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
