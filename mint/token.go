package mint

import (
	"github.com/kooksee/usmint/types"
	"github.com/ethereum/go-ethereum/common"
)

func NewToken(address common.Address) *types.Token {
	return types.NewToken(address, db)
}
