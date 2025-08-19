package middlewares

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"mini-alt/auth"
	"mini-alt/handlers/api"
	"net/http"
)

type S3ErrorResponse struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	Message   string   `xml:"Message"`
	RequestID string   `xml:"RequestId,omitempty"`
	HostID    string   `xml:"HostId,omitempty"`
}

func respondS3Error(c *gin.Context, code, message string) {
	c.XML(http.StatusUnauthorized, S3ErrorResponse{
		Code:    code,
		Message: message,
	})
	c.Abort()
}

func APIAuthenticationMiddleware(h *api.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, ok := c.Get("presignedAuth"); ok && val == true {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		dateHeader := c.GetHeader("x-amz-date")
		payloadHash := c.GetHeader("x-amz-content-sha256")

		if authHeader == "" || dateHeader == "" || payloadHash == "" {
			respondS3Error(c, "AccessDenied", "Missing Authentication Headers")
			return
		}

		parsedAuth, err := auth.ParseAuthorizationHeader(authHeader)
		if err != nil {
			respondS3Error(c, "AccessDenied", "Invalid Authorization Header")
			return
		}

		secretKey, err := h.Store.GetSecretKey(parsedAuth.AccessKeyID)
		if err != nil || secretKey == "" {
			respondS3Error(c, "AccessDenied", "Unknown Access Key")
			return
		}

		expectedSig, err := auth.CalculateSignature(c.Request, parsedAuth, secretKey, dateHeader, payloadHash)
		if err != nil {
			respondS3Error(c, "AccessDenied", err.Error())
			return
		}

		if parsedAuth.Signature != expectedSig {
			respondS3Error(c, "AccessDenied", "Signature Mismatch")
			return
		}

		c.Set("accessKey", parsedAuth.AccessKeyID)
		c.Next()
	}
}
