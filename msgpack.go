package encoding

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

type MsgPackCodec struct {
	codec
}

func (c *MsgPackCodec) Marshal(value interface{}) ([]byte, error) {

	if primitive, err := c.marshalPrimitive(value); err != errNotPrimitive {
		return primitive, err
	}

	var b bytes.Buffer

	enc := msgpack.GetEncoder()
	enc.Reset(&b)
	enc.UseCompactInts(true)

	err := enc.Encode(value)

	msgpack.PutEncoder(enc)

	if err != nil {
		return nil, err
	}

	return c.compress(b.Bytes())
}

func (c *MsgPackCodec) Unmarshal(byt []byte, ptr interface{}) (err error) {

	if err := c.unmarshalPrimitive(byt, ptr); err != errNotPrimitive {
		return err
	}

	if byt, err = c.decompress(byt); err != nil {
		return
	}

	return msgpack.Unmarshal(byt, ptr)
}

func NewMsgPackCodec(compressor Compressor) *MsgPackCodec {
	return &MsgPackCodec{
		codec{compressor},
	}
}
