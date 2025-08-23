package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	ExpiresAt string `json:"expiresAt"`
	Admin     bool   `json:"admin"`
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var request UserRegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessKey, _, err := h.Store.PutCredentials("", "", request.ExpiresAt, true, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newId, err := h.Store.RegisterUser(request.Username, request.Password, accessKey, request.ExpiresAt, request.Admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.Store.AddCredentialsOwner(accessKey, newId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
