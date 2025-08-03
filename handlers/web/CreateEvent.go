package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type EventRequest struct {
	Name   string `json:"name"`
	Bucket int64  `json:"bucket"`
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Store.CreateEvent(req.Name, req.Bucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}
