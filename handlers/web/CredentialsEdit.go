package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CredentialsEditRequest struct {
	AccessKey   string `json:"accessKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	ExpiresAt   string `json:"expiresAt"`
}

func (h *Handler) CredentialsEdit(c *gin.Context) {
	var request CredentialsEditRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Store.EditCredentials(request.AccessKey, request.Name, request.Description, request.ExpiresAt, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
