package yaml

import (
	"bytes"
	"io"
)

// Marshal serializes a Go value into YAML formatted bytes.
// Usage example:
//   data, err := yaml.Marshal(myStruct)
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal decodes a YAML byte slice into a target Go pointer.
// Usage example:
//   var val MyType
//   err := yaml.Unmarshal(data, &val)
func Unmarshal(in []byte, out interface{}) error {
	dec := NewDecoder(bytes.NewReader(in))
	return dec.Decode(out)
}

// NewDecoderStream constructs a Decoder targeting an io.Reader stream.
// Usage example:
//   dec := yaml.NewDecoderStream(r)
func NewDecoderStream(r io.Reader) *Decoder {
	return NewDecoder(r)
}

// NewEncoderStream constructs an Encoder targeting an io.Writer stream.
// Usage example:
//   enc := yaml.NewEncoderStream(w)
func NewEncoderStream(w io.Writer) *Encoder {
	return NewEncoder(w)
}
