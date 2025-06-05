package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

// DeleteObject receives the key of the file and removes that file.
// If no file is found, it does not return an error, it just returns.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (h *Handler) DeleteObject(c *gin.Context) {
	bucket := c.Param("bucket")
	objectKey := c.Param("object")

	h.Store.DeleteObject(bucket, objectKey)
	path := filepath.Join("uploads", bucket, objectKey)

	_ = os.Remove(path)

	c.Status(http.StatusNoContent)
}
