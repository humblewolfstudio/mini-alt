package web

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage/disk"
	"net/http"
)

func (h *Handler) GetSystemSpecs(c *gin.Context) {
	systemSpecs, err := disk.GetSystemSpecs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, systemSpecs)
}
