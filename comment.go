package yaml

// CommentPreserver associates header, line, and footer comments with AST nodes.
type CommentPreserver struct {
	comments map[int]string
}

func NewCommentPreserver() *CommentPreserver {
	return &CommentPreserver{comments: make(map[int]string)}
}
