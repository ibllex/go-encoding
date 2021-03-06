package encoding

import (
	"bytes"
	"encoding/gob"
)

type GobCodec struct {
	codec
}

func (c *GobCodec) Marshal(value interface{}) ([]byte, error) {

	if primitive, err := c.marshalPrimitive(value); err != errNotPrimitive {
		return primitive, err
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}

	return c.compress(b.Bytes())
}

func (c *GobCodec) Unmarshal(byt []byte, ptr interface{}) (err error) {

	if err := c.unmarshalPrimitive(byt, ptr); err != errNotPrimitive {
		return err
	}

	if byt, err = c.decompress(byt); err != nil {
		return
	}

	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)
	return decoder.Decode(ptr)
}

func NewGobCodec(compressor Compressor) *GobCodec {
	return &GobCodec{
		codec{compressor},
	}
}
