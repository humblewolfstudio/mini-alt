package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *WebHandler) AuthenticateUser(c *gin.Context) {
	id, err := c.Cookie("id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.Store.AuthenticateUser(idInt64, token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": idInt64})
}
