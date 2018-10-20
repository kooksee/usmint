package mint

import (
	"github.com/kooksee/usmint/mint/minter"
	"github.com/kooksee/usmint/mint/state"
	"github.com/kooksee/usmint/mint/validator"
	"github.com/tendermint/tendermint/libs/log"
)

func Init(logger log.Logger) {
	validator.Init(logger)
	minter.Init(logger)
	state.Init(logger)
}
