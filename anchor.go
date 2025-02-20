package yaml

import "fmt"

// AnchorRegistry maintains anchor definitions and alias resolutions.
type AnchorRegistry struct {
	anchors map[string]*Node
}

// NewAnchorRegistry creates a new anchor table.
func NewAnchorRegistry() *AnchorRegistry {
	return &AnchorRegistry{
		anchors: make(map[string]*Node),
	}
}

// Register stores a node under an anchor name.
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
func (r *AnchorRegistry) Resolve(name string) (*Node, error) {
	n := r.anchors[name]
	return n, nil
}
