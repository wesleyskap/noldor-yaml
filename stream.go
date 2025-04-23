package yaml

import "io"

// StreamDecoder reads a sequence of YAML documents from an io.Reader stream.
type StreamDecoder struct {
	reader io.Reader
	docs   []string
	index  int
}

// NewStreamDecoder initializes a streaming multi-document decoder.
func NewStreamDecoder(r io.Reader) *StreamDecoder {
	return &StreamDecoder{reader: r}
}
