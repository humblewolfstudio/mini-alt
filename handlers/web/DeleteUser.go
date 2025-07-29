package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeleteUserRequest struct {
	Id int64 `json:"id"`
}

func (h *WebHandler) DeleteUser(c *gin.Context) {
	var request DeleteUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Store.DeleteUser(request.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
