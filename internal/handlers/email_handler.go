package handlers

import (
	"net/http"
	"soundtube/internal/services"
	"soundtube/pkg"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	service *services.EmailService
	logger  *pkg.CustomLogger
}

func NewEmailHandler(service *services.EmailService, logger *pkg.CustomLogger) *EmailHandler {
	return &EmailHandler{service: service, logger: logger}
}

func (h *EmailHandler) VerifyEmail(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "EmailHandler.VerifyEmail")
	defer span.End()

	var req struct {
		Email       string `json:"emal"`
		VerifyToken string `json:"verify_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid params", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err := h.service.SendVerificationEmail(ctx, req.Email, req.VerifyToken)
	if err != nil {
		h.logger.Error("verify service error", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"messege": "user verified"})
}
