package mint

import (
	"github.com/kooksee/usmint/mint/minter"
	"github.com/kooksee/usmint/mint/state"
	"github.com/kooksee/usmint/mint/validator"
)

func Init() {
	validator.Init()
	minter.Init()
	state.Init()
}
