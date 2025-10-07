package middleware

import (
	"errors"
	"net/http"
	"soundtube/internal/services"
	"soundtube/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(s *services.LoginService, l *pkg.CustomLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")

		l.Info("Auth middleware started",
			"path", ctx.Request.URL.Path,
			"method", ctx.Request.Method,
			"authorization_header", tokenStr,
			"header_length", len(tokenStr),
		)

		if tokenStr == "" {
			var err = errors.New("emty token")
			l.Error("missing authorization header", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		username, err := s.ValidToken(ctx.Request.Context(), tokenStr)
		if err != nil {
			l.Error("invalid token", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("username", username)
		ctx.Set("token", tokenStr)

		l.Info("request authorized ", "username", username)
		ctx.Next()
	}
}
