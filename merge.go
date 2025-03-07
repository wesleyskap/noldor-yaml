package yaml

// MergeKeyResolver handles <<: *anchor map merging operations.
type MergeKeyResolver struct {
	registry *AnchorRegistry
}

// NewMergeKeyResolver constructs a resolver for merge keys.
func NewMergeKeyResolver(registry *AnchorRegistry) *MergeKeyResolver {
	return &MergeKeyResolver{
		registry: registry,
	}
}
