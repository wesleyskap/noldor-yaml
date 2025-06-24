package yaml

// ExplicitTagResolver resolves custom and standard YAML tag prefixes.
type ExplicitTagResolver struct {
	tags map[string]string
}

func NewExplicitTagResolver() *ExplicitTagResolver {
	return &ExplicitTagResolver{tags: make(map[string]string)}
}
