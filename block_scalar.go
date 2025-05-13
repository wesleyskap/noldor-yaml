package yaml

// BlockScalarFormatter formats block scalar styles.
type BlockScalarFormatter struct {
	indent int
	chomping string
}

func NewBlockScalarFormatter(indent int, chomping string) *BlockScalarFormatter {
	return &BlockScalarFormatter{indent: indent, chomping: chomping}
}
