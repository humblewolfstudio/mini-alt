package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/models"
	"mini-alt/storage/db"
)

type Handler struct {
	Store *db.Store
}

func GetUserFromContext(c *gin.Context) (models.User, bool) {
	u, exists := c.Get("user")
	if !exists {
		return models.User{}, false
	}
	user, ok := u.(models.User)
	return user, ok
}
