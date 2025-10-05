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
