package handlers

import (
	"net/http"
	"soundtube/internal/services"
	"soundtube/pkg"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

type LoginHandler struct {
	service *services.LoginService
	logger  *pkg.CustomLogger
}

func NewLoginHandler(service *services.LoginService, logger *pkg.CustomLogger) *LoginHandler {
	return &LoginHandler{service: service, logger: logger}
}

// Login authenticates user and returns JWT token
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string "Login successful"
// @Failure 400 {object} map[string]string "Invalid input format"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /api/auth/login [post]
func (h *LoginHandler) Login(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "LoginHandler.Login")
	defer span.End()

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(JsonInputFormat, err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	span.SetAttributes(
		attribute.String("user.name", req.Username),
	)

	token, err := h.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		h.logger.Error("login failed", err).WithTrace(ctx)
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout invalidates user token
// @Summary User logout
// @Description Invalidate user JWT token by adding to blacklist
// @Tags authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "JWT token to invalidate"
// @Success 200 {object} nil "Logout successful"
// @Failure 400 {object} map[string]string "Invalid input format"
// @Router /api/auth/logout [post]
func (h *LoginHandler) Logout(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "LoginHandler.Login")
	defer span.End()

	var req struct {
		Token string `json:"token"`
	}

	span.SetAttributes(
		attribute.String("token", req.Token),
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(JsonInputFormat, err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.Token); err != nil {
		h.logger.Error("logout failed", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
