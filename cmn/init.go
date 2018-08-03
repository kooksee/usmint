package cmn

import "github.com/kooksee/cmn"

var ErrPipe = cmn.Err.ErrWithMsg
var MustNotErr = cmn.Err.MustNotErr
var JsonMarshal = cmn.Json.Marshal
var JsonUnmarshal = cmn.Json.Unmarshal

type M map[string]interface{}
func (m M) String() string {
	d, _ := JsonMarshal(m)
	return string(d)
}
