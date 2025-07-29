package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeleteCredentialsRequest struct {
	AccessKey string `json:"accessKey"`
}

func (h *Handler) DeleteCredentials(c *gin.Context) {
	var request DeleteCredentialsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Store.DeleteCredentials(request.AccessKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credentials deleted"})
}
