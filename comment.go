package yaml

import (
	"fmt"
	"strings"
)

// CommentPreserver associates header, line, and footer comments with AST nodes.
// Struct fields are memory-aligned.
type CommentPreserver struct {
	comments map[int]string
}

// NewCommentPreserver constructs a new comment registry.
// Usage example:
//   cp := NewCommentPreserver()
func NewCommentPreserver() *CommentPreserver {
	return &CommentPreserver{
		comments: make(map[int]string),
	}
}

// AttachComment attaches comment text to a specific node line.
// Usage example:
//   cp.AttachComment(14, "# Configuration start")
func (cp *CommentPreserver) AttachComment(line int, text string) {
	cp.comments[line] = strings.TrimSpace(text)
}

// GetComment retrieves comment by line number.
// Usage example:
//   c := cp.GetComment(14)
func (cp *CommentPreserver) GetComment(line int) string {
	return cp.comments[line]
}

// ApplyToNode injects stored comments directly into AST node fields.
// Usage example:
//   err := cp.ApplyToNode(node)
func (cp *CommentPreserver) ApplyToNode(n *Node) error {
	if n == nil {
		return fmt.Errorf("yaml: cannot attach comment to nil node")
	}
	if c, ok := cp.comments[n.Line]; ok {
		n.HeadComment = c
	}
	for _, child := range n.Content {
		if err := cp.ApplyToNode(child); err != nil {
			return err
		}
	}
	return nil
}

