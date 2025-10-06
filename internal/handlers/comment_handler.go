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

func (h *CommentHandler) GetComments(c *gin.Context) {

}

func (h *CommentHandler) CreateComment(c *gin.Context) {

}

func (h *CommentHandler) UpdateComment(c *gin.Context) {

}

func (h *CommentHandler) DeleteComment(c *gin.Context) {

}
