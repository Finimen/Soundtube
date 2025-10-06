package middleware

import (
	"net/http"
	"soundtube/pkg"
	"soundtube/scripts"

	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware(r *pkg.RateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := scripts.GetClientIP(ctx.Request)

		if !r.Allow(ip) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": "too many requests"})
			return
		}

		ctx.Next()
	}
}
