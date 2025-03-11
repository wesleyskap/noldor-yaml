package yaml

import "fmt"

type MergeKeyResolver struct {
	registry *AnchorRegistry
}

func NewMergeKeyResolver(registry *AnchorRegistry) *MergeKeyResolver {
	return &MergeKeyResolver{registry: registry}
}

func (m *MergeKeyResolver) ApplyMerge(mapping *Node) error {
	if mapping == nil || mapping.Kind != MappingNode {
		return fmt.Errorf("yaml: merge key can only be applied to MappingNode, got %v", mapping)
	}
	for i := 0; i < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		if keyNode.Value == "<<" {
			// Found merge key token
		}
	}
	return nil
}
