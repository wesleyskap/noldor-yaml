package yaml

import (
	"fmt"
	"strings"
)

// AnchorRegistry maintains anchor definitions and alias resolutions.
// Struct fields are memory-aligned from largest to smallest byte size.
type AnchorRegistry struct {
	anchors map[string]*Node
}

// NewAnchorRegistry creates a new anchor table.
// Usage example:
//   reg := NewAnchorRegistry()
func NewAnchorRegistry() *AnchorRegistry {
	return &AnchorRegistry{
		anchors: make(map[string]*Node),
	}
}

// Register stores a node under an anchor name.
// Usage example:
//   err := reg.Register("base_config", node)
func (r *AnchorRegistry) Register(name string, n *Node) error {
	if name == "" {
		return fmt.Errorf("yaml: invalid empty anchor name, expected non-empty string identifier")
	}
	if n == nil {
		return fmt.Errorf("yaml: cannot register nil node for anchor %q", name)
	}
	r.anchors[name] = n
	return nil
}

// Resolve retrieves a referenced node by anchor name.
// Usage example:
//   node, err := reg.Resolve("base_config")
func (r *AnchorRegistry) Resolve(name string) (*Node, error) {
	n, exists := r.anchors[name]
	if !exists {
		return nil, fmt.Errorf("yaml: unknown anchor reference %q, expected registered anchor", name)
	}
	return n, nil
}

// ResolveAliases walks an AST node recursively and substitutes AliasNodes.
// Usage example:
//   err := reg.ResolveAliases(rootNode)
func (r *AnchorRegistry) ResolveAliases(n *Node) error {
	if n == nil {
		return nil
	}
	for i, child := range n.Content {
		if child.Kind == AliasNode && child.Value != "" {
			target, err := r.Resolve(strings.TrimPrefix(child.Value, "*"))
			if err != nil {
				return err
			}
			n.Content[i] = target
		} else {
			if err := r.ResolveAliases(child); err != nil {
				return err
			}
		}
	}
	return nil
}
