package middlewares

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func NormalizeObjectKeys() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := c.Params
		for i := range params {
			if params[i].Key == "object" {
				params[i].Value = strings.TrimPrefix(params[i].Value, "/")
			}
		}
		c.Params = params
		c.Next()
	}
}
