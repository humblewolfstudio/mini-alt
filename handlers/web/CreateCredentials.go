package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateCredentialsRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ExpiresAt   string `json:"expiresAt"`
}

func (h *Handler) CreateCredentials(c *gin.Context) {
	var request CreateCredentialsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessKey, secretKey, err := h.Store.CreateCredentials(request.Name, request.Description, request.ExpiresAt, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_key": accessKey, "secret_key": secretKey})
}
