package handlers

import (
	"errors"
	"net/http"
	"soundtube/internal/domain/sound"
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

// GetSounds retrieves all sounds
// @Summary Get all sounds
// @Description Get a list of all available sounds
// @Tags sounds
// @Security BearerAuth
// @Produce json
// @Success 200 {array} sound.SoundDTO "List of sounds"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/sounds [get]
func (h *SoundHandler) GetSounds(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "SoundHandler.GetSounds")
	defer span.End()
	sounds, err := h.service.GetSounds(ctx)
	if err != nil {
		h.logger.Error("get sound error", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	soundDTOs := sound.SoundsToDTO(sounds)

	h.logger.Info("Getted " + strconv.Itoa(len(soundDTOs)) + " from storage").WithTrace(ctx)
	c.JSON(http.StatusOK, soundDTOs)
}

// CreateSound creates a new sound record
// @Summary Create sound
// @Description Create a new sound record (without file upload)
// @Tags sounds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateSoundRequest true "Sound data"
// @Success 200 {object} string "Sound created successfully"
// @Failure 400 {object} map[string]string "Invalid input format"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/sounds [post]
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

// UpdateSound updates an existing sound
// @Summary Update sound
// @Description Update sound information
// @Tags sounds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Sound ID"
// @Param request body UpdateSoundRequest true "Update sound data"
// @Success 200 {object} object "Sound updated successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 403 {object} map[string]string "Forbidden - not sound owner"
// @Failure 404 {object} map[string]string "Sound not found"
// @Router /api/sounds/{id} [patch]
func (h *SoundHandler) UpdateSound(c *gin.Context) {
	//TODO: Implement update func
}

// DeleteSound deletes a sound
// @Summary Delete sound
// @Description Delete a sound by name
// @Tags sounds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Sound ID"
// @Param request body DeleteSoundRequest true "Delete data"
// @Success 200 {object} object "Sound deleted successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 403 {object} map[string]string "Forbidden - not sound owner"
// @Failure 404 {object} map[string]string "Sound not found"
// @Router /api/sounds/{id} [delete]
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
