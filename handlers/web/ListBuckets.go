package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListBuckets(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	buckets, err := h.Store.ListBuckets(id.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, buckets)
}
