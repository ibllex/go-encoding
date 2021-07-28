package encoding_test

import (
	"reflect"
	"testing"

	"github.com/ibllex/go-encoding"
)

type MockStruct struct {
	X int
}

func (s MockStruct) Method1() {}

type Interface1 interface {
	Method1()
}

var (
	structType                 = MockStruct{1}
	ptrStruct                  = &MockStruct{2}
	emptyIface     interface{} = MockStruct{3}
	iface1         Interface1  = MockStruct{4}
	sliceStruct                = []MockStruct{{5}, {6}, {7}}
	ptrSliceStruct             = []*MockStruct{{8}, {9}, {10}}

	valueMap = map[string]interface{}{
		"bytes":          []byte{0x61, 0x62, 0x63, 0x64},
		"string":         "string",
		"bool":           true,
		"int":            5,
		"int8":           int8(5),
		"int16":          int16(5),
		"int32":          int32(5),
		"int64":          int64(5),
		"uint":           uint(5),
		"uint8":          uint8(5),
		"uint16":         uint16(5),
		"uint32":         uint32(5),
		"uint64":         uint64(5),
		"float32":        float32(5),
		"float64":        float64(5),
		"array":          [5]int{1, 2, 3, 4, 5},
		"slice":          []int{1, 2, 3, 4, 5},
		"emptyIf":        emptyIface,
		"Iface1":         iface1,
		"map":            map[string]string{"foo": "bar"},
		"ptrStruct":      ptrStruct,
		"structType":     structType,
		"sliceStruct":    sliceStruct,
		"ptrSliceStruct": ptrSliceStruct,
	}
)

func TestRoundTrip(t *testing.T) {
	compressor := encoding.NewS2Compressor(4)

	codecs := map[string]encoding.Codec{
		"gob":     encoding.NewGobCodec(compressor),
		"msgpack": encoding.NewMsgPackCodec(compressor),
	}

	for name, codec := range codecs {
		for _, expected := range valueMap {
			bytes, err := codec.Marshal(expected)
			if err != nil {
				t.Error(err)
				continue
			}

			ptrActual := reflect.New(reflect.TypeOf(expected)).Interface()
			err = codec.Unmarshal(bytes, ptrActual)
			if err != nil {
				t.Error(err)
				continue
			}

			actual := reflect.ValueOf(ptrActual).Elem().Interface()
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("[%s] (expected) %T %v != %T %v (actual)", name, expected, expected, actual, actual)
			}
		}

	}
}
