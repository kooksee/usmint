package types

import "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Map map[string]interface{}

func (m Map) String() string {
	d, _ := json.MarshalToString(m)
	return d
}
