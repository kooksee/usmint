package kts

import (
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"time"
	"github.com/ethereum/go-ethereum/common"
	"encoding/hex"
	"github.com/kooksee/usmint/cmn"
)

func DecodeContract(data string) (*Contract, error) {
	dt, err := hex.DecodeString(data)
	if err != nil {
		return nil, cmn.ErrPipe("DecodeContract 1", err)
	}

	c := &Contract{}
	if err := cmn.ErrPipe("DecodeContract 2", cmn.JsonUnmarshal(dt, c)); err != nil {
		return nil, err
	}

	return c, nil
}

type Contract struct {
	Address []byte `json:"addr,omitempty"`
	Method  string `json:"method,omitempty"`
	Data    []byte `json:"data,omitempty"`
}

// 得到合约的地址
func CreateContractAddress(addr []byte, code []byte) common.Address {
	return common.BytesToAddress(crypto.Keccak256(addr, code, big.NewInt(time.Now().Unix()).Bytes()))
}

func DecodeSCParams(data string) {
}
