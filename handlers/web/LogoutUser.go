package web

import "github.com/gin-gonic/gin"

func (h *Handler) LogoutUser(c *gin.Context) {
	c.SetCookie("id", "", -1, "/", "", false, false)
	c.SetCookie("token", "", -1, "/", "", false, false)
	c.SetCookie("admin", "false", -1, "/", "", false, false)

	c.JSON(200, gin.H{"message": "ok"})
}
