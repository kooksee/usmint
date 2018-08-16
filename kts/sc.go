package kts

import "github.com/ethereum/go-ethereum/crypto"

type Contract struct {
	Address         []byte `json:"addr,omitempty"`
	Method          string `json:"method,omitempty"`
	Data            []byte `json:"data,omitempty"`
	ContractAddress []byte
	Code            []byte
	Tx              []byte
}

func (c *Contract) CreateContractAddress(tx []byte) []byte {
	return crypto.Keccak256(tx, c.Code)
}

func DecodeSCParams(data string) {
}
