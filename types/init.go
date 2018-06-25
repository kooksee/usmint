package types

import "github.com/json-iterator/go"
import cmn "github.com/tendermint/tmlibs/common"

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Address in go-crypto style
type Address = cmn.HexBytes

func DataResult(data interface{}) string {
	d, _ := json.MarshalToString(map[string]interface{}{"data": data})
	return d
}

type Map map[string]interface{}

func (m Map) String() string {
	d, _ := json.MarshalToString(m)
	return d
}
