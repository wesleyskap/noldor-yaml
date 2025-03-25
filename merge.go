package yaml

import (
	"fmt"
	"strings"
)

// MergeKeyResolver handles '<<: *anchor' map merging operations.
// Struct fields are memory-aligned.
type MergeKeyResolver struct {
	registry *AnchorRegistry
}

// NewMergeKeyResolver constructs a resolver for merge keys.
// Usage example:
//   m := NewMergeKeyResolver(reg)
func NewMergeKeyResolver(registry *AnchorRegistry) *MergeKeyResolver {
	return &MergeKeyResolver{
		registry: registry,
	}
}

// ApplyMerge merges anchor node fields into target mapping nodes.
// Usage example:
//   err := m.ApplyMerge(mappingNode)
func (m *MergeKeyResolver) ApplyMerge(mapping *Node) error {
	if mapping == nil || mapping.Kind != MappingNode {
		return fmt.Errorf("yaml: merge key can only be applied to MappingNode, got %v", mapping)
	}
	var newContent []*Node
	for i := 0; i < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		valNode := mapping.Content[i+1]
		if keyNode.Value == "<<" {
			resolvedVal := valNode
			if valNode.Kind == AliasNode && m.registry != nil {
				target, err := m.registry.Resolve(strings.TrimPrefix(valNode.Value, "*"))
				if err != nil {
					return err
				}
				resolvedVal = target
			}
			if resolvedVal.Kind == MappingNode {
				newContent = append(newContent, resolvedVal.Content...)
			}
		} else {
			newContent = append(newContent, keyNode, valNode)
		}
	}
	mapping.Content = newContent
	return nil
}

// Graceful nil registry handling
