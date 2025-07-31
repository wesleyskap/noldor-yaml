package yaml

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Parser manages tokenizing and AST construction for YAML input streams.
type Parser struct {
	reader *bufio.Reader
	line   int
}

// NewParser constructs a Parser bound to an input reader.
// Usage example:
//   p := NewParser(r)
func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
		line:   1,
	}
}

// ParseTree parses the input stream into a root DocumentNode AST.
// Usage example:
//   doc, err := p.ParseTree()
func (p *Parser) ParseTree() (*Node, error) {
	data, err := io.ReadAll(p.reader)
	if err != nil {
		return nil, fmt.Errorf("yaml: failed reading input stream: %w", err)
	}
	doc := &Node{
		Kind: DocumentNode,
		Line: 1,
	}
	lines := strings.Split(string(data), "\n")
	root, err := parseLines(lines, 0, 0)
	if err != nil {
		return nil, err
	}
	if root != nil {
		doc.Content = append(doc.Content, root)
	}
	return doc, nil
}

// parseLines converts a set of indentation-based lines into a Node AST hierarchy.
func parseLines(lines []string, index int, minIndent int) (*Node, error) {
	lines = stripEmptyAndComments(lines)
	if len(lines) == 0 {
		return &Node{Kind: ScalarNode, Value: ""}, nil
	}
	first := lines[0]
	trimmed := strings.TrimSpace(first)
	if strings.HasPrefix(trimmed, "- ") || trimmed == "-" {
		return parseSequenceLines(lines)
	}
	if isMappingLine(trimmed) {
		return parseMappingLines(lines)
	}
	return parseScalarLine(trimmed), nil
}

// stripEmptyAndComments filters blank lines and extract line comments.
func stripEmptyAndComments(lines []string) []string {
	var result []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		result = append(result, l)
	}
	return result
}

// isMappingLine checks if a line contains a key: value pair format.
func isMappingLine(line string) bool {
	idx := strings.Index(line, ":")
	if idx <= 0 {
		return false
	}
	if idx == len(line)-1 {
		return true
	}
	nextChar := line[idx+1]
	return nextChar == ' ' || nextChar == '\t' || nextChar == '\n'
}

// parseScalarLine produces a scalar Node from a single line string.
func parseScalarLine(trimmed string) *Node {
	style := Style(0)
	val := trimmed
	if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) && len(val) >= 2 {
		style = DoubleQuotedStyle
		val = val[1 : len(val)-1]
	} else if strings.HasPrefix(val, `'`) && strings.HasSuffix(val, `'`) && len(val) >= 2 {
		style = SingleQuotedStyle
		val = val[1 : len(val)-1]
	}
	return &Node{
		Kind:  ScalarNode,
		Value: val,
		Style: style,
	}
}

// parseSequenceLines parses sequential array block elements starting with '- '.
func parseSequenceLines(lines []string) (*Node, error) {
	seq := &Node{
		Kind: SequenceNode,
	}
	items := splitSequenceItems(lines)
	for _, itemLines := range items {
		child, err := parseSequenceChild(itemLines)
		if err != nil {
			return nil, err
		}
		seq.Content = append(seq.Content, child)
	}
	return seq, nil
}

// parseSequenceChild converts a block item into an array element node.
func parseSequenceChild(itemLines []string) (*Node, error) {
	if len(itemLines) == 0 {
		return &Node{Kind: ScalarNode, Value: ""}, nil
	}
	first := strings.TrimSpace(itemLines[0])
	if strings.HasPrefix(first, "- ") {
		first = strings.TrimPrefix(first, "- ")
	} else if first == "-" {
		first = ""
	}
	itemLines[0] = first
	return parseLines(itemLines, 0, 0)
}

// splitSequenceItems partitions sequence lines by bullet boundaries.
func splitSequenceItems(lines []string) [][]string {
	var items [][]string
	var current []string
	baseIndent := getIndent(lines[0])
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		indent := getIndent(l)
		if (trimmed == "-" || strings.HasPrefix(trimmed, "- ")) && indent == baseIndent {
			if len(current) > 0 {
				items = append(items, current)
			}
			current = []string{l}
		} else {
			current = append(current, l)
		}
	}
	if len(current) > 0 {
		items = append(items, current)
	}
	return items
}

// parseMappingLines converts map block key-value lines into a MappingNode.
func parseMappingLines(lines []string) (*Node, error) {
	mapping := &Node{
		Kind: MappingNode,
	}
	entries := splitMappingEntries(lines)
	for key, valLines := range entries {
		keyNode := &Node{Kind: ScalarNode, Value: key}
		valNode, err := parseLines(valLines, 0, 0)
		if err != nil {
			return nil, err
		}
		mapping.Content = append(mapping.Content, keyNode, valNode)
	}
	return mapping, nil
}

// splitMappingEntries groups mapping lines by top-level keys.
func splitMappingEntries(lines []string) map[string][]string {
	entries := make(map[string][]string)
	var currentKey string
	var currentVal []string
	baseIndent := getIndent(lines[0])
	for _, line := range lines {
		indent := getIndent(line)
		trimmed := strings.TrimSpace(line)
		if indent == baseIndent && isMappingLine(trimmed) {
			if currentKey != "" {
				entries[currentKey] = currentVal
			}
			k, v := splitKeyValue(trimmed)
			currentKey = k
			currentVal = []string{}
			if v != "" {
				currentVal = append(currentVal, v)
			}
		} else {
			currentVal = append(currentVal, line)
		}
	}
	if currentKey != "" {
		entries[currentKey] = currentVal
	}
	return entries
}

// splitKeyValue extracts key and value from a mapping line.
func splitKeyValue(line string) (string, string) {
	idx := strings.Index(line, ":")
	key := strings.TrimSpace(line[:idx])
	val := strings.TrimSpace(line[idx+1:])
	return key, val
}

// getIndent returns leading whitespace count for a line.
func getIndent(line string) int {
	count := 0
	for _, ch := range line {
		if ch == ' ' {
			count++
		} else if ch == '\t' {
			count += 2
		} else {
			break
		}
	}
	return count
}

// ParseBytes helper function to parse raw bytes into a Node tree.
// Usage example:
//   node, err := ParseBytes([]byte("key: val"))
func ParseBytes(b []byte) (*Node, error) {
	p := NewParser(bytes.NewReader(b))
	return p.ParseTree()
}

