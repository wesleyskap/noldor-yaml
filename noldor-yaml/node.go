// Package yaml provides YAML encoding, decoding, and AST node parsing.
package yaml

import (
	"fmt"
)

// Kind represents the structural type of a YAML Node.
type Kind uint32

const (
	// DocumentNode represents a top-level document node.
	DocumentNode Kind = 1 << iota
	// SequenceNode represents an array or sequence node.
	SequenceNode
	// MappingNode represents a key-value mapping node.
	MappingNode
	// ScalarNode represents a primitive value node.
	ScalarNode
	// AliasNode represents a reference node.
	AliasNode
)

// Style represents the formatting style of a YAML scalar or collection.
type Style uint32

const (
	// TaggedStyle preserves explicit tag annotations.
	TaggedStyle Style = 1 << iota
	// DoubleQuotedStyle formats scalars with double quotes.
	DoubleQuotedStyle
	// SingleQuotedStyle formats scalars with single quotes.
	SingleQuotedStyle
	// LiteralStyle formats scalars as literal block text.
	LiteralStyle
	// FoldedStyle formats scalars as folded block text.
	FoldedStyle
	// FlowStyle formats collections using inline bracket syntax.
	FlowStyle
)

// Node represents an individual element in a YAML syntax tree.
// Struct fields are aligned from largest to smallest byte size to optimize memory layout.
type Node struct {
	Content     []*Node
	HeadComment string
	LineComment string
	FootComment string
	Tag         string
	Value       string
	Alias       *Node
	Line        int
	Column      int
	Kind        Kind
	Style       Style
}

// Decode unmarshals the node tree into the target Go value.
// Usage example:
//   var data map[string]string
//   err := node.Decode(&data)
func (n *Node) Decode(target interface{}) error {
	if n == nil {
		return fmt.Errorf("yaml: cannot decode nil node into %T, expected valid node reference", target)
	}
	decoder := NewDecoderFromNode(n)
	return decoder.Decode(target)
}

// Encode populates the node tree from a Go value.
// Usage example:
//   var node yaml.Node
//   err := node.Encode(map[string]string{"key": "value"})
func (n *Node) Encode(value interface{}) error {
	encoder := NewNodeEncoder(n)
	return encoder.Encode(value)
}

// IsZero checks whether the node is uninitialized.
// Usage example:
//   if node.IsZero() { ... }
func (n *Node) IsZero() bool {
	if n == nil {
		return true
	}
	return n.Kind == 0 && n.Value == "" && len(n.Content) == 0
}
