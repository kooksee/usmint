package consts

import "strings"

const (
	// events
	ValidatorSet = "ValidatorSet"

	// db
	Metadata = "metadata"
	License  = "license"

	// prefix
	DbPrefix        = "db:"
	AccountPrefix   = "act:"
	ValidatorPrefix = "val:"

	TokenPrefix = "token:"

	MetaPrefix  = "meta:"
	StatePrefix = "state"
	MinerPrefix = "miner"

	TokenAddress = "0x0000000000000000000000000"
	TotalSupply  = "totalSupply"
)

func Meta(ms ... string) string {
	return MetaPrefix + strings.Join(ms, "")
}
