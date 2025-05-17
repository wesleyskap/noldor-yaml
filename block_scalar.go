package yaml

import (
	"fmt"
	"strings"
)

// BlockScalarFormatter formats block scalar styles (Literal | and Folded >).
// Struct fields are memory-aligned.
type BlockScalarFormatter struct {
	indent int
	chomping string
}

// NewBlockScalarFormatter initializes block scalar handler.
// Usage example:
//   b := NewBlockScalarFormatter(2, "strip")
func NewBlockScalarFormatter(indent int, chomping string) *BlockScalarFormatter {
	return &BlockScalarFormatter{
		indent:   indent,
		chomping: chomping,
	}
}

// FormatLiteral formats multiline strings using literal style '|'.
// Usage example:
//   out := b.FormatLiteral("line1\nline2")
func (b *BlockScalarFormatter) FormatLiteral(text string) string {
	lines := strings.Split(text, "\n")
	pad := strings.Repeat(" ", b.indent)
	header := "|\n"
	if b.chomping == "strip" {
		header = "|-\n"
	}
	var sb strings.Builder
	sb.WriteString(header)
	for i, l := range lines {
		if i == len(lines)-1 && l == "" {
			continue
		}
		sb.WriteString(pad + l + "\n")
	}
	return sb.String()
}

// ParseLiteralBlock extracts text from literal block lines.
// Usage example:
//   txt, err := b.ParseLiteralBlock(lines)
func (b *BlockScalarFormatter) ParseLiteralBlock(lines []string) (string, error) {
	if len(lines) == 0 {
		return "", fmt.Errorf("yaml: empty block scalar lines, expected content lines")
	}
	var clean []string
	for _, l := range lines {
		clean = append(clean, strings.TrimPrefix(l, strings.Repeat(" ", b.indent)))
	}
	return strings.Join(clean, "\n"), nil
}

// Strip chomping indicator
