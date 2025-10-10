package handlers

import (
	"errors"
	"net/http"
	"soundtube/internal/services"
	"soundtube/pkg"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReactionHandler struct {
	service *services.ReactionService
	logger  *pkg.CustomLogger
}

func NewReactionHandler(service *services.ReactionService, logger *pkg.CustomLogger) *ReactionHandler {
	return &ReactionHandler{service: service, logger: logger}
}

func (h *ReactionHandler) SetReactionSound(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "ReactionHandler.SetReactionSound")
	defer span.End()

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("invalid user_id in context", errors.New("user_id not found")).WithTrace(ctx)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		h.logger.Error("invalid user_id type", errors.New("type assertion failed")).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	soundIDStr := c.Param("id")
	soundID, err := strconv.Atoi(soundIDStr)
	if err != nil {
		h.logger.Error("invalid soind id", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		ReactionType string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ReactionType != "like" && req.ReactionType != "dislike" {
		h.logger.Error("invalid reaction type", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, gin.H{"error": "reaction type must be like or dislike"})
		return
	}

	err = h.service.SetSoundReaction(ctx, userID, soundID, req.ReactionType)
}

func (h *ReactionHandler) DeleteReactionSound(c *gin.Context) {
	//TODO
}

func (h *ReactionHandler) GetReactionSound(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "ReactionHandler.GetReactionSound")
	defer span.End()

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("invalid user_id in context", errors.New("user_id not found"))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		h.logger.Error("invalid user_id type", errors.New("type assertion failed"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	soundIDStr := c.Param("id")
	soundID, err := strconv.Atoi(soundIDStr)
	if err != nil {
		h.logger.Error("invalid sound id", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sound ID"})
		return
	}

	reactions, err := h.service.GetSoundReactions(ctx, userID, soundID)
	if err != nil {
		h.logger.Error("failed to get reactions", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reactions"})
		return
	}

	c.JSON(http.StatusOK, reactions)
}

func (h *ReactionHandler) SetReactionComment(c *gin.Context) {
	//TODO
}

func (h *ReactionHandler) DeleteReactionComment(c *gin.Context) {
	//TODO
}

func (h *ReactionHandler) GetReactionComment(c *gin.Context) {
	//TODO
}
