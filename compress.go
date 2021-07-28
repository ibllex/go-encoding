package encoding

import (
	"bytes"
	"fmt"

	"github.com/klauspost/compress/s2"
)

const (
	S2Compression = 0x1
)

type DefaultCompressor struct {
	way   int
	chunk int
}

func (c *DefaultCompressor) Compress(data []byte) ([]byte, error) {

	switch c.way {
	case S2Compression:
		n := s2.MaxEncodedLen(len(data)) + 1
		b := make([]byte, n, n+c.chunk)
		b = s2.Encode(b, data)
		b = append(b, S2Compression)
		return b, nil
	default:
		return data, fmt.Errorf("unknown compression method: %x", c.way)
	}

}

func (c *DefaultCompressor) Decompress(data []byte) ([]byte, error) {

	switch c := data[len(data)-1]; c {
	case S2Compression:
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

func NewS2Compressor(chunk int) Compressor {
	return &DefaultCompressor{
		way:   S2Compression,
		chunk: chunk,
	}
}
