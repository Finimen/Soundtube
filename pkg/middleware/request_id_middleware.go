package middleware

import (
	"soundtube/scripts"

	"github.com/gin-gonic/gin"
)

func RequsetIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requsetID := ctx.GetHeader("X-Request-ID")
		if requsetID == "" {
			requsetID = scripts.GenerateUUID()
		}

		ctx.Set("request_id", requsetID)

		ctx.Header("X-Request-ID", requsetID)

		ctx.Next()
	}
}
