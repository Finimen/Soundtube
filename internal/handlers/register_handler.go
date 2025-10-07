package handlers

import (
	"net/http"
	"soundtube/internal/services"
	"soundtube/pkg"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

type RegisterHandler struct {
	service *services.RegisterService
	logger  *pkg.CustomLogger
}

func NewRegisterHandler(service *services.RegisterService, logger *pkg.CustomLogger) *RegisterHandler {
	return &RegisterHandler{service: service, logger: logger}
}

func (h *RegisterHandler) Register(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "RegisterHandler.Register")
	defer span.End()

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(JsonInputFormat, err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info(req.Username, req.Password, req.Email)

	span.SetAttributes(
		attribute.String("user.username", req.Username),
		attribute.String("user.password", req.Password),
		attribute.String("user.email", req.Email),
	)

	err := h.service.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		h.logger.Error("register service error", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
