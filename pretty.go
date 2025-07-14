package yaml

// PrettyPrinter generates human-readable indented YAML representations.
type PrettyPrinter struct {
	indent int
}

func NewPrettyPrinter(indent int) *PrettyPrinter {
	if indent <= 0 {
		indent = 2
	}
	return &PrettyPrinter{indent: indent}
}
