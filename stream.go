package yaml

import "io"

// StreamDecoder reads a sequence of YAML documents from an io.Reader stream.
type StreamDecoder struct {
	reader io.Reader
	docs   []string
	index  int
}
