package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage"
	"mini-alt/utils"
	"net/http"
)

// DeleteBucket reveices the bucket name and removes that bucket with all of its content
// If no bucket is found (or the bucket name is incorrect) it throws an error.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html
func (h *ApiHandler) DeleteBucket(c *gin.Context) {
	bucketName := c.Param("bucket")
	if err := storage.DeleteBucket(bucketName); err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError", "Could not delete bucket", bucketName)
		return
	}

	err := h.Store.DeleteBucket(bucketName)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError", "Could not delete bucket", bucketName)
		return
	}

	c.Status(http.StatusNoContent)
}
