package yaml

import (
	"io"
	"strings"
)

type StreamDecoder struct {
	reader io.Reader
	docs   []string
	index  int
}

func NewStreamDecoder(r io.Reader) *StreamDecoder {
	data, _ := io.ReadAll(r)
	parts := strings.Split(string(data), "\n---")
	return &StreamDecoder{reader: r, docs: parts, index: 0}
}
