package encoding

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

type MsgPackCodec struct {
	compressor Compressor
}

func (c *MsgPackCodec) Marshal(value interface{}) ([]byte, error) {
	switch value := value.(type) {
	case nil:
		return nil, nil
	case []byte:
		return value, nil
	case string:
		return []byte(value), nil
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

	if c.compressor != nil {
		return c.compressor.Compress(b.Bytes())
	}

	return b.Bytes(), nil
}

func (c *MsgPackCodec) Unmarshal(b []byte, value interface{}) error {
	if len(b) == 0 {
		return nil
	}

	switch value := value.(type) {
	case nil:
		return nil
	case *[]byte:
		clone := make([]byte, len(b))
		copy(clone, b)
		*value = clone
		return nil
	case *string:
		*value = string(b)
		return nil
	}

	if c.compressor != nil {
		var err error
		b, err = c.compressor.Decompress(b)
		if err != nil {
			return err
		}
	}

	return msgpack.Unmarshal(b, value)
}

func NewMsgPackCodec(compressor Compressor) *MsgPackCodec {
	return &MsgPackCodec{
		compressor: compressor,
	}
}
