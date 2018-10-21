package minter

import (
	"testing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kooksee/usmint/cmn/wire"
	"github.com/kooksee/usmint/kts"
	"fmt"
)

func TestName(t *testing.T) {
	m := &SetMiner{Addr: common.BytesToAddress([]byte("oooo")), Power: 2}

	var h kts.DataHandler
	if err := wire.Decode(m.Encode(), &h); err != nil {
		panic(err.Error())
	}
	fmt.Println(h)
}
