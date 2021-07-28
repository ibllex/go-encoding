package encoding

import (
	"bytes"
	"encoding/json"
)

type JsonCodec struct {
	codec
}

func (c *JsonCodec) Marshal(value interface{}) ([]byte, error) {

	if primitive, err := c.marshalPrimitive(value); err != errNotPrimitive {
		return primitive, err
	}

	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}

	return c.compress(b.Bytes())
}

func (c *JsonCodec) Unmarshal(byt []byte, ptr interface{}) (err error) {

	if err := c.unmarshalPrimitive(byt, ptr); err != errNotPrimitive {
		return err
	}

	if byt, err = c.decompress(byt); err != nil {
		return
	}

	b := bytes.NewBuffer(byt)
	decoder := json.NewDecoder(b)
	return decoder.Decode(ptr)
}

func NewJsonCodec(compressor Compressor) *JsonCodec {
	return &JsonCodec{
		codec{compressor},
	}
}
