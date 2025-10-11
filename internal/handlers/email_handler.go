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

// VerifyEmail verifies user email with token
// @Summary Verify email
// @Description Verify user email address using verification token
// @Tags authentication
// @Produce json
// @Param token query string true "Email verification token"
// @Success 200 {object} map[string]string "Email verified successfully"
// @Failure 400 {object} map[string]string "Token is required or invalid"
// @Router /api/auth/verify-email [get]
func (h *EmailHandler) VerifyEmail(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "EmailHandler.VerifyEmail")
	defer span.End()

	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	if err := h.service.VerifyEmail(ctx, token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
