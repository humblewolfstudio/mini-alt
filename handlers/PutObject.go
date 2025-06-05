package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage"
	"mini-alt/utils"
	"net/http"
)

// PutObject receives the bucket name, the object key and the object and persists it.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (h *Handler) PutObject(c *gin.Context, bucket, object string) {
	path, err := storage.CreateObjectFilePath(bucket, object)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create directory path.", bucket)
		return
	}

	written, err := storage.CreateObject(path, c.Request.Body)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create file.", bucket)
		return
	}

	h.Store.PutObject(bucket, object, written)
	c.Status(http.StatusOK)
}
