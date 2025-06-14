package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage"
	"net/http"
)

// DeleteObject receives the key of the file and removes that file.
// If no file is found, it does not return an error, it just returns.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (h *Handler) DeleteObject(c *gin.Context) {
	bucketName := c.Param("bucketName")
	objectKey := c.Param("object")

	// Also delete all files
	err := h.Store.DeleteObject(bucketName, objectKey)
	if err == nil {
		storage.DeleteObjectFile(bucketName, objectKey)
	}

	c.Status(http.StatusNoContent)
}
