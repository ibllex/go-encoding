package encoding

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/s2"
)

const (
	s2Compression   = 0x1
	gzipCompression = 0x2
)

type DefaultCompressor struct {
	algorithm int
}

//
// gzip compression algorithm
//

func (c *DefaultCompressor) gzipCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	w.Write(data)
	err := w.Close()
	return b.Bytes(), err
}

func (c *DefaultCompressor) gzipDecompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return data, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

//
// s2 compression algorithm
//

func (c *DefaultCompressor) s2Compress(data []byte) ([]byte, error) {
	n := s2.MaxEncodedLen(len(data))
	b := make([]byte, n)
	b = s2.Encode(b, data)
	return b, nil
}

func (c *DefaultCompressor) s2Decompress(data []byte) ([]byte, error) {
	_, err := s2.DecodedLen(data)
	if err != nil {
		return data, err
	}

	var buf bytes.Buffer
	return s2.Decode(buf.Bytes(), data)
}

//
// Compressor interface
//

func (c *DefaultCompressor) Compress(data []byte) ([]byte, error) {

	switch c.algorithm {
	case s2Compression:
		return c.s2Compress(data)
	case gzipCompression:
		return c.gzipCompress(data)
	default:
		return data, fmt.Errorf("unknown compression algorithm: %x", c.algorithm)
	}

}

func (c *DefaultCompressor) Decompress(data []byte) ([]byte, error) {

	switch c.algorithm {
	case s2Compression:
		return c.s2Decompress(data)
	case gzipCompression:
		return c.gzipDecompress(data)
	default:
		return data, fmt.Errorf("unknown compression algorithm: %x", c.algorithm)
	}

}

func NewGzipCompressor() Compressor {
	return &DefaultCompressor{
		algorithm: gzipCompression,
	}
}

func NewS2Compressor() Compressor {
	return &DefaultCompressor{
		algorithm: s2Compression,
	}
}
