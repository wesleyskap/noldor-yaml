package yaml

import (
	"fmt"
	"io"
	"strings"
)

// StreamDecoder reads a sequence of YAML documents from an io.Reader stream.
// Struct fields are memory-aligned.
type StreamDecoder struct {
	reader io.Reader
	docs   []string
	index  int
}

// NewStreamDecoder initializes a streaming multi-document decoder.
// Usage example:
//   dec := NewStreamDecoder(reader)
func NewStreamDecoder(r io.Reader) *StreamDecoder {
	data, _ := io.ReadAll(r)
	parts := strings.Split(string(data), "\n---")
	var cleaned []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" && trimmed != "..." {
			cleaned = append(cleaned, trimmed)
		}
	}
	return &StreamDecoder{
		reader: r,
		docs:   cleaned,
		index:  0,
	}
}

// More checks if there are remaining documents in the stream.
// Usage example:
//   for dec.More() { ... }
func (d *StreamDecoder) More() bool {
	return d.index < len(d.docs)
}

// Decode decodes the next document into target.
// Usage example:
//   err := dec.Decode(&config)
func (d *StreamDecoder) Decode(target interface{}) error {
	if !d.More() {
		return fmt.Errorf("yaml: EOF reached, no more documents in stream")
	}
	content := d.docs[d.index]
	d.index++
	return Unmarshal([]byte(content), target)
}
