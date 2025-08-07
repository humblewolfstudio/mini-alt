package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/storage/disk"
	"mini-alt/utils"
	"net/http"
)

// GetObject gets the bucket and the file key and returns the file.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (h *Handler) GetObject(c *gin.Context, bucket string, objectKey string) {
	path, err := disk.GetObjectPath(bucket, objectKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchKey", "Object not found.", bucket)
		return
	}

	go events.HandleEventObject(h.Store, types.EventGet, utils.ClearObjectKeyWithBucket(bucket, objectKey), bucket, "")

	c.File(path)
}
