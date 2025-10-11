package handlers

import (
	"soundtube/pkg"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	logger *pkg.CustomLogger
}

func NewCommentHandler(logger *pkg.CustomLogger) *CommentHandler {
	return &CommentHandler{logger: logger}
}

// GetComments retrieves comments for a specific sound
// @Summary Get sound comments
// @Description Get all comments for a specific sound
// @Tags comments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Sound ID"
// @Success 200 {array} object "List of comments"
// @Failure 400 {object} map[string]string "Invalid sound ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/sounds/{id}/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	//TODO
}

// CreateComment creates a new comment for a sound
// @Summary Create comment
// @Description Add a new comment to a sound
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Sound ID"
// @Param request body CreateCommentRequest true "Comment data"
// @Success 201 {object} object "Comment created successfully"
// @Failure 400 {object} map[string]string "Invalid input or sound ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/sounds/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	//TODO
}

// UpdateComment updates an existing comment
// @Summary Update comment
// @Description Update a comment by ID
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param request body UpdateCommentRequest true "Updated comment data"
// @Success 200 {object} object "Comment updated successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 403 {object} map[string]string "Forbidden - not comment owner"
// @Failure 404 {object} map[string]string "Comment not found"
// @Router /api/comments/{id} [patch]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	//TODO
}

// DeleteComment deletes a comment
// @Summary Delete comment
// @Description Delete a comment by ID
// @Tags comments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} object "Comment deleted successfully"
// @Failure 403 {object} map[string]string "Forbidden - not comment owner"
// @Failure 404 {object} map[string]string "Comment not found"
// @Router /api/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	//TODO
}
