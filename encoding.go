package encoding

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
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}
