package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage"
	"mini-alt/utils"
	"net/http"
)

func (h *Handler) DeleteBucket(c *gin.Context) {
	bucketName := c.Param("bucket")
	if err := storage.DeleteBucket(bucketName); err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError", "Could not delete bucket", bucketName)
		return
	}

	h.Store.DeleteBucket(bucketName)

	c.Status(http.StatusNoContent)
}
