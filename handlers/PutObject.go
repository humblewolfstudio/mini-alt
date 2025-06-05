package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"mini-alt/utils"
	"net/http"
	"os"
	"path/filepath"
)

// PutObject receives the bucket name, the object key and the object and persists it.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (h *Handler) PutObject(c *gin.Context, bucket, object string) {
	dst := filepath.Join("uploads", bucket, object)
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create directory path.", bucket)
		return
	}

	file, err := os.Create(dst)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create file.", bucket)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	written, err := io.Copy(file, c.Request.Body)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not write file.", bucket)
		return
	}

	h.Store.PutObject(bucket, object, written)
	c.Status(http.StatusOK)
}
