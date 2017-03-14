package yaml

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Decoder reads and decodes YAML values from an input stream.
// Struct fields are aligned from largest to smallest byte size to optimize memory layout.
type Decoder struct {
	reader io.Reader
	node   *Node
	line   int
}

// NewDecoder creates a Decoder reading from an io.Reader.
// Usage example:
//   dec := yaml.NewDecoder(r)
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: r,
		line:   1,
	}
}

// NewDecoderFromNode creates a Decoder targeting an existing AST Node.
// Usage example:
//   dec := yaml.NewDecoderFromNode(node)
func NewDecoderFromNode(node *Node) *Decoder {
	return &Decoder{
		node: node,
		line: node.Line,
	}
}

// Decode unmarshals the next YAML document into target pointer v.
// Usage example:
//   err := dec.Decode(&myStruct)
func (d *Decoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("yaml: invalid decode target %T, expected non-nil pointer", v)
	}
	node, err := d.fetchRootNode()
	if err != nil {
		return err
	}
	return unmarshalNode(node, rv.Elem())
}

// fetchRootNode reads input data into AST Node if node is not pre-set.
func (d *Decoder) fetchRootNode() (*Node, error) {
	if d.node != nil {
		return d.node, nil
	}
	p := NewParser(d.reader)
	root, err := p.ParseTree()
	if err != nil {
		return nil, err
	}
	if root.Kind == DocumentNode && len(root.Content) > 0 {
		return root.Content[0], nil
	}
	return root, nil
}

// unmarshalNode dispatches AST unmarshaling based on Go reflect Kind.
func unmarshalNode(node *Node, target reflect.Value) error {
	if target.Type() == reflect.TypeOf(Node{}) {
		target.Set(reflect.ValueOf(*node))
		return nil
	}
	if target.Type() == reflect.TypeOf(&Node{}) {
		target.Set(reflect.ValueOf(node))
		return nil
	}
	switch target.Kind() {
	case reflect.Bool:
		return decodeBoolNode(node, target)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return decodeIntNode(node, target)
	case reflect.Float32, reflect.Float64:
		return decodeFloatNode(node, target)
	case reflect.String:
		target.SetString(node.Value)
		return nil
	case reflect.Slice:
		return decodeSliceNode(node, target)
	case reflect.Map:
		return decodeMapNode(node, target)
	case reflect.Struct:
		return decodeStructNode(node, target)
	case reflect.Interface:
		return decodeInterfaceNode(node, target)
	default:
		return fmt.Errorf("yaml: unsupported target kind %s for scalar value %q", target.Kind(), node.Value)
	}
}

// decodeBoolNode converts scalar string to boolean value.
func decodeBoolNode(node *Node, target reflect.Value) error {
	b, err := ParseBool(node.Value)
	if err != nil {
		return fmt.Errorf("yaml: cannot unmarshal scalar %q into bool, expected boolean value (true, false, yes, no)", node.Value)
	}
	target.SetBool(b)
	return nil
}

// decodeIntNode converts scalar string to integer value.
func decodeIntNode(node *Node, target reflect.Value) error {
	n, err := ParseInt(node.Value)
	if err != nil {
		return fmt.Errorf("yaml: cannot unmarshal scalar %q into int type, expected valid integer shape: %w", node.Value, err)
	}
	target.SetInt(n)
	return nil
}

// decodeFloatNode converts scalar string to floating-point value.
func decodeFloatNode(node *Node, target reflect.Value) error {
	f, err := ParseFloat(node.Value)
	if err != nil {
		return fmt.Errorf("yaml: cannot unmarshal scalar %q into float type, expected float shape: %w", node.Value, err)
	}
	target.SetFloat(f)
	return nil
}

// decodeSliceNode populates a Go slice from a SequenceNode.
func decodeSliceNode(node *Node, target reflect.Value) error {
	if node.Kind != SequenceNode {
		return fmt.Errorf("yaml: cannot unmarshal node of kind %d into slice, expected sequence format", node.Kind)
	}
	sliceType := target.Type()
	elemType := sliceType.Elem()
	slice := reflect.MakeSlice(sliceType, 0, len(node.Content))
	for _, child := range node.Content {
		elemPtr := reflect.New(elemType)
		if err := unmarshalNode(child, elemPtr.Elem()); err != nil {
			return err
		}
		slice = reflect.Append(slice, elemPtr.Elem())
	}
	target.Set(slice)
	return nil
}

// decodeMapNode populates a Go map from a MappingNode.
func decodeMapNode(node *Node, target reflect.Value) error {
	if node.Kind != MappingNode {
		return fmt.Errorf("yaml: cannot unmarshal node of kind %d into map, expected mapping format", node.Kind)
	}
	mapType := target.Type()
	if target.IsNil() {
		target.Set(reflect.MakeMap(mapType))
	}
	keyType := mapType.Key()
	valType := mapType.Elem()
	for i := 0; i < len(node.Content); i += 2 {
		kNode := node.Content[i]
		vNode := node.Content[i+1]
		kVal := reflect.New(keyType).Elem()
		vVal := reflect.New(valType).Elem()
		if err := unmarshalNode(kNode, kVal); err != nil {
			return err
		}
		if err := unmarshalNode(vNode, vVal); err != nil {
			return err
		}
		target.SetMapIndex(kVal, vVal)
	}
	return nil
}

// decodeStructNode maps MappingNode fields to struct exported fields.
func decodeStructNode(node *Node, target reflect.Value) error {
	if node.Kind != MappingNode {
		return fmt.Errorf("yaml: cannot unmarshal node of kind %d into struct %s, expected mapping format", node.Kind, target.Type())
	}
	fieldMap := buildStructFieldMap(target)
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		valNode := node.Content[i+1]
		if fieldVal, found := fieldMap[strings.ToLower(key)]; found {
			if err := unmarshalNode(valNode, fieldVal); err != nil {
				return err
			}
		}
	}
	return nil
}

// buildStructFieldMap indexes struct fields by yaml tag or field name.
func buildStructFieldMap(target reflect.Value) map[string]reflect.Value {
	res := make(map[string]reflect.Value)
	t := target.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tag := field.Tag.Get("yaml")
		name := field.Name
		if tag != "" && tag != "-" {
			parts := strings.Split(tag, ",")
			name = parts[0]
		}
		res[strings.ToLower(name)] = target.Field(i)
	}
	return res
}

// decodeInterfaceNode infers dynamic type for interface{} target.
func decodeInterfaceNode(node *Node, target reflect.Value) error {
	switch node.Kind {
	case ScalarNode:
		target.Set(reflect.ValueOf(node.Value))
	case SequenceNode:
		var slice []interface{}
		sv := reflect.ValueOf(&slice).Elem()
		if err := decodeSliceNode(node, sv); err != nil {
			return err
		}
		target.Set(sv)
	case MappingNode:
		m := make(map[string]interface{})
		mv := reflect.ValueOf(m)
		if err := decodeMapNode(node, mv); err != nil {
			return err
		}
		target.Set(mv)
	}
	return nil
}
