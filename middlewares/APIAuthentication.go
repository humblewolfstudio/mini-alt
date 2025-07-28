package middlewares

import (
	"github.com/gin-gonic/gin"
	"mini-alt/auth"
	"mini-alt/handlers/api"
	"net/http"
)

func APIAuthenticationMiddleware(h *api.ApiHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		dateHeader := c.GetHeader("x-amz-date")
		payloadHash := c.GetHeader("x-amz-content-sha256")

		if authHeader == "" || dateHeader == "" || payloadHash == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing auth headers"})
			return
		}

		parsedAuth, err := auth.ParseAuthorizationHeader(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth header"})
			return
		}

		secretKey, err := h.Store.GetSecretKey(parsedAuth.AccessKeyID)
		if secretKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unknown access key"})
			return
		}

		expectedSig := auth.CalculateSignature(c.Request, parsedAuth, secretKey, dateHeader, payloadHash)

		if parsedAuth.Signature != expectedSig {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Signature mismatch"})
			return
		}

		c.Next()
	}
}
