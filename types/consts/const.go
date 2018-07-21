package consts

import "strings"

const (
	// events
	DbSet        = "DbSet"
	DbGet        = "DbGet"
	ValidatorSet = "ValidatorSet"
	AccountSet   = "AccountSet"

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
