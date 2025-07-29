package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *WebHandler) ListUsers(c *gin.Context) {
	users, err := h.Store.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
