package luas

import (
	"github.com/yuin/gopher-lua"
)

func JsonDecode(L *lua.LState) int {
	value, err := decodeRaw(L, []byte(L.CheckString(1)))
	if err != nil {
		L.Push(lua.LNil)
		//logger.Error("json decode error", "err", err)
		return 1
	}
	L.Push(value)
	return 1
}

// Decode converts the JSON encoded data to Lua values.
func decodeRaw(L *lua.LState, data []byte) (lua.LValue, error) {
	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	return decode(L, value), nil
}

func decode(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case int64:
		return lua.LNumber(converted)
	case int:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(decode(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), decode(L, item))
		}
		return tbl
	}
	return lua.LNil
}
