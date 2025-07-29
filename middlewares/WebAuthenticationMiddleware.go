package middlewares

import (
	"github.com/gin-gonic/gin"
	"mini-alt/handlers/web"
	"net/http"
	"strconv"
)

func WebAuthenticationMiddleware(h *web.WebHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := c.Cookie("id")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = h.Store.AuthenticateUser(idInt64, token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("id", idInt64)

		c.Next()
	}
}
