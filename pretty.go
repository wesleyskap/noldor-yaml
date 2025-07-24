package yaml

import (
	"fmt"
	"strings"
)

// PrettyPrinter generates human-readable indented YAML representations.
// Struct fields are memory-aligned.
type PrettyPrinter struct {
	indent int
}

// NewPrettyPrinter constructs a new pretty printer with specified indentation.
// Usage example:
//   pp := NewPrettyPrinter(2)
func NewPrettyPrinter(indent int) *PrettyPrinter {
	if indent <= 0 {
		indent = 2
	}
	return &PrettyPrinter{indent: indent}
}

// PrintNode converts a Node AST hierarchy into a formatted YAML string.
// Usage example:
//   out, err := pp.PrintNode(docNode)
func (pp *PrettyPrinter) PrintNode(n *Node) (string, error) {
	if n == nil {
		return "", fmt.Errorf("yaml: cannot pretty print nil node")
	}
	var sb strings.Builder
	pp.writeNode(&sb, n, 0)
	return sb.String(), nil
}

func (pp *PrettyPrinter) writeNode(sb *strings.Builder, n *Node, level int) {
	pad := strings.Repeat(" ", level*pp.indent)
	switch n.Kind {
	case DocumentNode:
		for _, child := range n.Content {
			pp.writeNode(sb, child, level)
		}
	case MappingNode:
		for i := 0; i < len(n.Content); i += 2 {
			k := n.Content[i]
			v := n.Content[i+1]
			if v.Kind == MappingNode || v.Kind == SequenceNode {
				sb.WriteString(fmt.Sprintf("%s%s:\n", pad, k.Value))
				pp.writeNode(sb, v, level+1)
			} else {
				sb.WriteString(fmt.Sprintf("%s%s: %s\n", pad, k.Value, v.Value))
			}
		}
	case SequenceNode:
		for _, item := range n.Content {
			sb.WriteString(fmt.Sprintf("%s- %s\n", pad, item.Value))
		}
	case ScalarNode:
		sb.WriteString(fmt.Sprintf("%s%s\n", pad, n.Value))
	}
}

// MappingNode indentation support
// SequenceNode bullet rendering
// ScalarNode formatting
// Nil pointer validation
