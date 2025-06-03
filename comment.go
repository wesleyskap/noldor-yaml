package yaml

// CommentPreserver associates header, line, and footer comments with AST nodes.
type CommentPreserver struct {
	comments map[int]string
}
