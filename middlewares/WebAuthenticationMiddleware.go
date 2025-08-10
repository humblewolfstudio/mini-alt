package middlewares

import (
	"github.com/gin-gonic/gin"
	"mini-alt/handlers/web"
	"net/http"
	"strconv"
)

func WebAuthenticationMiddleware(h *web.Handler) gin.HandlerFunc {
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

func WebAuthenticationAdminMiddleware(h *web.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := c.Cookie("id")
		if err != nil {
			println("1")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := c.Cookie("token")
		if err != nil {
			println("2")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		admin, err := c.Cookie("admin")
		if err != nil || admin != "true" {
			println("3")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			println("4")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		println(idInt64)
		err = h.Store.AuthenticateAdmin(idInt64, token)
		if err != nil {
			println(err.Error())
			println("5")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("id", idInt64)

		c.Next()
	}
}
