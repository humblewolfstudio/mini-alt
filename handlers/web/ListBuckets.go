package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListBuckets(c *gin.Context) {
	buckets, err := h.Store.ListBuckets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, buckets)
}
