package handlers

import (
	"net/http"
	"soundtube/internal/services"
	"soundtube/pkg"

	"github.com/gin-gonic/gin"
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

	h.logger.Info("Getted ", len(sounds), " from storage").WithTrace(ctx)
	c.JSON(http.StatusOK, nil)
}

func (h *SoundHandler) CreateSound(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "SoundHandler.CreateSound")
	defer span.End()

	var req struct {
		Name     string `json:"name"`
		Album    string `json:"album"`
		Genre    string `json:"genre"`
		AuthorId int    `json"author_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(JsonInputFormat, err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err := h.service.CreateSound(ctx, req.Name, req.Album, req.Genre, req.AuthorId)
	if err != nil {
		h.logger.Error("get sound error", err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "req.Name"+" was created!")
}

func (h *SoundHandler) UpdateSound(c *gin.Context) {

}

func (h *SoundHandler) DeleteSound(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "SoundHandler.CreateSound")
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

	c.JSON(http.StatusOK, "req.Name"+" was deleted!")
}
