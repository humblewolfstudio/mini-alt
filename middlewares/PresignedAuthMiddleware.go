package middlewares

import (
	"github.com/gin-gonic/gin"
	"mini-alt/auth"
	"mini-alt/handlers/api"
	"time"
)

func PresignedAuthMiddleware(h *api.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("X-Amz-Algorithm") != "AWS4-HMAC-SHA256" {
			c.Next()
			return
		}

		credential := c.Query("X-Amz-Credential")
		signature := c.Query("X-Amz-Signature")
		date := c.Query("X-Amz-Date")
		expires := c.Query("X-Amz-Expires")
		signedHeaders := c.Query("X-Amz-SignedHeaders")

		if credential == "" || signature == "" || date == "" || expires == "" || signedHeaders == "" {
			respondS3Error(c, "AccessDenied", "Missing required query parameters")
			return
		}

		parsed, err := auth.ParseCredentialQuery(credential)
		if err != nil {
			respondS3Error(c, "AccessDenied", "Invalid credential format")
			return
		}
		parsed.SignedHeaders = signedHeaders
		parsed.Signature = signature

		expire, err := time.ParseDuration(expires + "s")
		if err != nil {
			respondS3Error(c, "AccessDenied", "Invalid X-Amz-Expires")
			return
		}
		startTime, err := time.Parse("20060102T150405Z", date)
		if err != nil {
			respondS3Error(c, "AccessDenied", "Invalid X-Amz-Date")
			return
		}

		if time.Now().UTC().After(startTime.Add(expire)) {
			respondS3Error(c, "AccessDenied", "Request has expired")
			return
		}

		secretKey, err := h.Store.GetSecretKey(parsed.AccessKeyID)
		if err != nil || secretKey == "" {
			respondS3Error(c, "AccessDenied", "Invalid Access Key")
			return
		}

		expectedSig, err := auth.CalculateSignaturePresigned(c.Request, parsed, secretKey)
		if err != nil {
			respondS3Error(c, "AccessDenied", "Signature calculation failed: "+err.Error())
			return
		}

		if !auth.SecureCompare(signature, expectedSig) {
			respondS3Error(c, "AccessDenied", "Signature mismatch")
			return
		}

		c.Set("presignedAuth", true)
		c.Next()
	}
}
