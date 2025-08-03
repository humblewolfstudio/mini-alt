package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListEvents(c *gin.Context) {
	events, err := h.Store.ListEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
