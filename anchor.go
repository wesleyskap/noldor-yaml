package yaml

// AnchorRegistry maintains anchor definitions and alias resolutions.
type AnchorRegistry struct {
	anchors map[string]*Node
}
