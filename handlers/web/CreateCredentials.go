package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *WebHandler) CreateCredentials(c *gin.Context) {
	accessKey, secretKey, err := h.Store.CreateCredentials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_key": accessKey, "secret_key": secretKey})
}
