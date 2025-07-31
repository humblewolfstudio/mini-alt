package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage/disk"
	"mini-alt/utils"
	"net/http"
)

// GetObject gets the bucket and the file key and returns the file.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (h *Handler) GetObject(c *gin.Context, bucketName string, objectKey string) {
	path, err := disk.GetObjectPath(bucketName, objectKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchKey", "Object not found.", bucketName)
		return
	}

	c.File(path)
}
