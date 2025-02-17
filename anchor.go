package yaml

// AnchorRegistry maintains anchor definitions and alias resolutions.
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
func (r *AnchorRegistry) Register(name string, n *Node) error {
	r.anchors[name] = n
	return nil
}
