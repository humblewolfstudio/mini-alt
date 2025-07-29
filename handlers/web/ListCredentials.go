package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *WebHandler) ListCredentials(c *gin.Context) {
	credentials, err := h.Store.ListCredentials()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, credentials)
}
