package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/utils"
	"net/http"
	"os"
	"path/filepath"
)

// GetObject gets the bucket and the file key and returns the file.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (h *Handler) GetObject(c *gin.Context, bucket string, object string) {
	path := filepath.Join("uploads", bucket, object)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchKey", "Object not found.", bucket)
		return
	}

	c.File(path)
}
