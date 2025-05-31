package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"mini-alt/utils"
	"net/http"
	"os"
	"path/filepath"
)

func (h *Handler) PutObjectOrBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")

	if object == "/" || object == "" {
		h.handleCreateBucket(c, bucket)
		return
	}

	object = object[1:]
	h.handleUploadObject(c, bucket, object)
}

func (h *Handler) handleUploadObject(c *gin.Context, bucket, object string) {
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
	defer file.Close()

	written, err := io.Copy(file, c.Request.Body)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not write file.", bucket)
		return
	}

	h.Store.PutObject(bucket, object, written)
	c.Status(http.StatusOK)
}

func (h *Handler) GetObject(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")

	path := filepath.Join("uploads", bucket, object)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		utils.RespondS3Error(c, http.StatusConflict, "ObjectNotFound",
			"Object not found.", bucket)
		return
	}

	c.File(path)
}
