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
)

func Meta(strs ... string) string {
	return MetaPrefix + strings.Join(strs, "")
}
