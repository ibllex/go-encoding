package encoding

import (
	"bytes"
	"fmt"

	"github.com/klauspost/compress/s2"
)

const (
	s2Compression = 0x1
)

type DefaultCompressor struct {
	way int
}

func (c *DefaultCompressor) Compress(data []byte) ([]byte, error) {

	switch c.way {
	case s2Compression:
		n := s2.MaxEncodedLen(len(data)) + 1
		b := make([]byte, n)
		b = s2.Encode(b, data)
		b = append(b, s2Compression)
		return b, nil
	default:
		return data, fmt.Errorf("unknown compression method: %x", c.way)
	}

}

func (c *DefaultCompressor) Decompress(data []byte) ([]byte, error) {

	switch c := data[len(data)-1]; c {
	case s2Compression:
		data = data[:len(data)-1]

		_, err := s2.DecodedLen(data)
		if err != nil {
			return data, err
		}

		var buf bytes.Buffer

		data, err = s2.Decode(buf.Bytes(), data)
		if err != nil {
			return data, err
		}
	default:
		return data, fmt.Errorf("unknown compression method: %x", c)
	}

	return data, nil
}

func NewS2Compressor() Compressor {
	return &DefaultCompressor{
		way: s2Compression,
	}
}
