package yaml

import (
	"fmt"
	"strings"
)

// ExplicitTagResolver resolves custom and standard YAML tag prefixes.
// Struct fields are memory-aligned.
type ExplicitTagResolver struct {
	tags map[string]string
}

// NewExplicitTagResolver initializes a tag handler.
// Usage example:
//   tr := NewExplicitTagResolver()
func NewExplicitTagResolver() *ExplicitTagResolver {
	return &ExplicitTagResolver{
		tags: make(map[string]string),
	}
}

// RegisterTag maps custom tag prefix to full URI.
// Usage example:
//   tr.RegisterTag("!custom", "tag:yaml.org,2002:str")
func (tr *ExplicitTagResolver) RegisterTag(handle, uri string) {
	tr.tags[handle] = uri
}

// ResolveTag expands shorthand tag handles into resolved tag strings.
// Usage example:
//   resolved := tr.ResolveTag("!custom")
func (tr *ExplicitTagResolver) ResolveTag(handle string) string {
	if uri, ok := tr.tags[handle]; ok {
		return uri
	}
	return handle
}

// ApplyTagToNode assigns an explicit tag to a target AST Node.
// Usage example:
//   err := tr.ApplyTagToNode(node, "!!binary")
func (tr *ExplicitTagResolver) ApplyTagToNode(n *Node, tag string) error {
	if n == nil {
		return fmt.Errorf("yaml: cannot assign tag %q to nil node", tag)
	}
	n.Tag = strings.TrimSpace(tag)
	n.Style |= TaggedStyle
	return nil
}

// Docstring with usage example
