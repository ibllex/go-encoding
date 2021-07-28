package encoding

import (
	"errors"
	"reflect"
	"strconv"
)

var errNotPrimitive = errors.New("value is not primitive")

// Marshaler transforms the given value into bytes
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
}

// Unmarshaler transforms bytes produced by Marshal back into a Go object
type Unmarshaler interface {
	Unmarshal([]byte, interface{}) error
}

// Codec is a combination of Unmarshaler and Marshaler
type Codec interface {
	Marshaler
	Unmarshaler
}

// Compressor compress and decompress the given data
type Compressor interface {
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
}

type codec struct {
	compressor Compressor
}

func (c *codec) marshalPrimitive(value interface{}) ([]byte, error) {

	if data, ok := value.([]byte); ok {
		return data, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}

	return nil, errNotPrimitive
}

func (c *codec) unmarshalPrimitive(byt []byte, ptr interface{}) error {
	if data, ok := ptr.(*[]byte); ok {
		*data = byt
		return nil
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr {
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
				return err
			}

			p.SetInt(i)
			return nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.ParseUint(string(byt), 10, 64)
			if err != nil {
				return err
			}

			p.SetUint(i)
			return nil
		}
	}

	return errNotPrimitive
}

func (c *codec) compress(data []byte) ([]byte, error) {
	if c.compressor == nil {
		return data, nil
	}

	return c.compressor.Compress(data)
}

func (c *codec) decompress(data []byte) ([]byte, error) {
	if c.compressor == nil {
		return data, nil
	}

	return c.compressor.Decompress(data)
}
