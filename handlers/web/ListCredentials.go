package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListCredentials(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	credentials, err := h.Store.ListCredentials(id.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, credentials)
}
