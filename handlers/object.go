package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"mini-alt/utils"
	"net/http"
	"os"
	"path/filepath"
)

// PutObjectOrBucket receives the endpoint of creating an object or a bucket (due to gin problem with * endpoints).
func (h *Handler) PutObjectOrBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")

	if object == "/" || object == "" {
		h.HandleCreateBucket(c, bucket)
		return
	}

	object = object[1:]
	h.HandlePutObject(c, bucket, object)
}

// HandlePutObject receives the bucket name, the object key and the object and persists it.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (h *Handler) HandlePutObject(c *gin.Context, bucket, object string) {
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

// GetObject gets the bucket and the file key and returns the file.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (h *Handler) GetObject(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")

	path := filepath.Join("uploads", bucket, object)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchKey",
			"Object not found.", bucket)
		return
	}

	c.File(path)
}
