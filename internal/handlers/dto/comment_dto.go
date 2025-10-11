// comment_dto.go
package handlers

// CreateCommentRequest represents the request body for creating a comment
type CreateCommentRequest struct {
	Content string `json:"content" example:"Great sound!"`
}

// UpdateCommentRequest represents the request body for updating a comment
type UpdateCommentRequest struct {
	Content string `json:"content" example:"Updated comment content"`
}
