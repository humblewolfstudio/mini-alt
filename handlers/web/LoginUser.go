package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *WebHandler) LoginUser(c *gin.Context) {
	var request UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Store.LoginUser(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("id", strconv.FormatInt(res.Id, 10), 3600, "/", "", false, true)
	c.SetCookie("token", res.Token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
