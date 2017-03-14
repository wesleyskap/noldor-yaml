package yaml

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

// Encoder serializes Go values into formatted YAML text streams.
// Struct fields are aligned from largest to smallest byte size to optimize memory layout.
type Encoder struct {
	writer io.Writer
	indent int
}

// NewEncoder initializes an Encoder bound to an io.Writer.
// Usage example:
//   enc := yaml.NewEncoder(w)
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: w,
		indent: 0,
	}
}

// NewNodeEncoder initializes an Encoder that populates a Node tree.
// Usage example:
//   enc := yaml.NewNodeEncoder(node)
func NewNodeEncoder(n *Node) *NodeEncoderHelper {
	return &NodeEncoderHelper{target: n}
}

// NodeEncoderHelper assists in encoding directly into a Node.
type NodeEncoderHelper struct {
	target *Node
}

// Encode encodes value into the target Node.
func (h *NodeEncoderHelper) Encode(v interface{}) error {
	node, err := valueToNode(reflect.ValueOf(v))
	if err != nil {
		return err
	}
	*h.target = *node
	return nil
}

// Encode converts value v into YAML output written to the destination stream.
// Usage example:
//   err := enc.Encode(myStruct)
func (e *Encoder) Encode(v interface{}) error {
	rv := reflect.ValueOf(v)
	node, err := valueToNode(rv)
	if err != nil {
		return err
	}
	text := renderNode(node, 0)
	_, err = io.WriteString(e.writer, text)
	return err
}

// valueToNode converts any reflect.Value into a Node AST.
func valueToNode(rv reflect.Value) (*Node, error) {
	if !rv.IsValid() {
		return &Node{Kind: ScalarNode, Value: "null"}, nil
	}
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return &Node{Kind: ScalarNode, Value: "null"}, nil
		}
		rv = rv.Elem()
	}
	if rv.Type() == reflect.TypeOf(Node{}) {
		n := rv.Interface().(Node)
		return &n, nil
	}
	switch rv.Kind() {
	case reflect.Struct:
		return structToNode(rv)
	case reflect.Map:
		return mapToNode(rv)
	case reflect.Slice, reflect.Array:
		return sliceToNode(rv)
	default:
		return scalarToNode(rv)
	}
}

// scalarToNode converts primitive types into scalar Nodes.
func scalarToNode(rv reflect.Value) (*Node, error) {
	valStr := fmt.Sprintf("%v", rv.Interface())
	return &Node{
		Kind:  ScalarNode,
		Value: valStr,
	}, nil
}

// sliceToNode converts slices/arrays into SequenceNodes.
func sliceToNode(rv reflect.Value) (*Node, error) {
	seq := &Node{Kind: SequenceNode}
	for i := 0; i < rv.Len(); i++ {
		child, err := valueToNode(rv.Index(i))
		if err != nil {
			return nil, err
		}
		seq.Content = append(seq.Content, child)
	}
	return seq, nil
}

// mapToNode converts maps into MappingNodes.
func mapToNode(rv reflect.Value) (*Node, error) {
	mapping := &Node{Kind: MappingNode}
	keys := rv.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprintf("%v", keys[i].Interface()) < fmt.Sprintf("%v", keys[j].Interface())
	})
	for _, k := range keys {
		keyNode, _ := valueToNode(k)
		valNode, err := valueToNode(rv.MapIndex(k))
		if err != nil {
			return nil, err
		}
		mapping.Content = append(mapping.Content, keyNode, valNode)
	}
	return mapping, nil
}

// structToNode converts struct fields into MappingNodes.
func structToNode(rv reflect.Value) (*Node, error) {
	mapping := &Node{Kind: MappingNode}
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tag := field.Tag.Get("yaml")
		if tag == "-" {
			continue
		}
		name := extractFieldName(field, tag)
		valNode, err := valueToNode(rv.Field(i))
		if err != nil {
			return nil, err
		}
		keyNode := &Node{Kind: ScalarNode, Value: name}
		mapping.Content = append(mapping.Content, keyNode, valNode)
	}
	return mapping, nil
}

// extractFieldName resolves field name from yaml tag or fallback name.
func extractFieldName(field reflect.StructField, tag string) string {
	if tag != "" {
		parts := strings.Split(tag, ",")
		if parts[0] != "" {
			return parts[0]
		}
	}
	return field.Name
}

// renderNode formats a Node tree into indented YAML string output.
func renderNode(n *Node, indent int) string {
	if n == nil {
		return ""
	}
	pad := strings.Repeat("  ", indent)
	switch n.Kind {
	case ScalarNode:
		return n.Value + "\n"
	case SequenceNode:
		return renderSequenceNode(n, indent, pad)
	case MappingNode:
		return renderMappingNode(n, indent, pad)
	default:
		if len(n.Content) > 0 {
			return renderNode(n.Content[0], indent)
		}
		return ""
	}
}

// renderSequenceNode formats array elements with '-' bullet prefixes.
func renderSequenceNode(n *Node, indent int, pad string) string {
	var sb strings.Builder
	for _, child := range n.Content {
		sb.WriteString(pad + "- " + renderNode(child, indent+1))
	}
	return sb.String()
}

// renderMappingNode formats key-value mapping elements.
func renderMappingNode(n *Node, indent int, pad string) string {
	var sb strings.Builder
	for i := 0; i < len(n.Content); i += 2 {
		key := n.Content[i].Value
		valNode := n.Content[i+1]
		if valNode.Kind == ScalarNode {
			sb.WriteString(fmt.Sprintf("%s%s: %s", pad, key, renderNode(valNode, indent)))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s:\n%s", pad, key, renderNode(valNode, indent+1)))
		}
	}
	return sb.String()
}
