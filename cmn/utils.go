package cmn

import (
	"github.com/tendermint/tmlibs/common"
	"math"
	"encoding/binary"
)

var Fmt = common.Fmt

func BytesTrimSpace(bs []byte) []byte {
	for i, b := range bs {
		if b != 0 {
			return bs[i : len(bs)-1]
		}
	}
	return nil
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

// 使用二进制存储整形
func Int64ToByte(x int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(x))
	return b
}

func ByteToInt64(x []byte) int64 {
	return int64(binary.BigEndian.Uint64(x))
}

func SortStruct(s interface{}) ([]byte, error) {
	s1, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	b := map[string]interface{}{}
	if err := json.Unmarshal(s1, &b); err != nil {
		return nil, err
	}
	b1, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return b1, nil
}

func Fn(fn func(args ... interface{}) error) error {
	return fn()
}
