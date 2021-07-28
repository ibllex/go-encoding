package encoding

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"strconv"
)

type GobCodec struct {
	compressor Compressor
}

func (c *GobCodec) Marshal(value interface{}) ([]byte, error) {

	if data, ok := value.([]byte); ok {
		return data, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}

	if c.compressor != nil {
		return c.compressor.Compress(b.Bytes())
	}

	return b.Bytes(), nil
}

func (c *GobCodec) Unmarshal(byt []byte, ptr interface{}) (err error) {
	if data, ok := ptr.(*[]byte); ok {
		*data = byt
		return
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr {
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			i, err = strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
			} else {
				p.SetInt(i)
			}
			return

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var i uint64
			i, err = strconv.ParseUint(string(byt), 10, 64)
			if err != nil {
			} else {
				p.SetUint(i)
			}
			return
		}
	}

	if c.compressor != nil {
		byt, err = c.compressor.Decompress(byt)
		if err != nil {
			return
		}
	}

	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)
	if err = decoder.Decode(ptr); err != nil {
		return
	}

	return
}

func NewGobCodec(compressor Compressor) *GobCodec {
	return &GobCodec{
		compressor: compressor,
	}
}
