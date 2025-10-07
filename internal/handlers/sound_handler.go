package handlers

import (
	"errors"
	"net/http"
	"soundtube/internal/services"
	"soundtube/pkg"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

type SoundHandler struct {
	service *services.SoundService
	logger  *pkg.CustomLogger
}

func NewSoundHandler(service *services.SoundService, logger *pkg.CustomLogger) *SoundHandler {
	return &SoundHandler{service: service, logger: logger}
}

func (h *SoundHandler) GetSounds(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "SoundHandler.GetSounds")
	defer span.End()
	sounds, err := h.service.GetSounds(ctx)
	if err != nil {
		h.logger.Error("get sound error", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	h.logger.Info("Getted " + strconv.Itoa(len(sounds)) + " from storage").WithTrace(ctx)
	c.JSON(http.StatusOK, sounds)
}

func (h *SoundHandler) CreateSound(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "SoundHandler.CreateSound")
	defer span.End()

	h.logger.Info("Context keys:").WithTrace(ctx)
	for _, key := range c.Keys {
		val, _ := c.Get(key)
		h.logger.Info("Key: %s, Value: %v, Type: %T", key, val, val).WithTrace(ctx)
	}

	userIDValid, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("user_id not found in context", errors.New("user not found")).WithTrace(ctx)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, ok := userIDValid.(int)
	if !ok {
		h.logger.Error("user_id is not integer", errors.New("invalid type")).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var req struct {
		Name  string `json:"name"`
		Album string `json:"album"`
		Genre string `json:"genre"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(JsonInputFormat, err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	span.SetAttributes(
		attribute.String("sound.name", req.Name),
		attribute.String("sound.album", req.Album),
		attribute.String("sound.genre", req.Genre),
	)

	err := h.service.CreateSound(ctx, req.Name, req.Album, req.Genre, userID)
	if err != nil {
		h.logger.Error("get sound error", err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "req.Name"+" was created!")
}

func (h *SoundHandler) UpdateSound(c *gin.Context) {
	//TODO: Implement update func
}

func (h *SoundHandler) DeleteSound(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "SoundHandler.DeleteSound")
	defer span.End()

	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(JsonInputFormat, err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err := h.service.DeleteSound(ctx, req.Name)
	if err != nil {
		h.logger.Error("get sound error", err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"messege": req.Name + " was deleted!"})
}
