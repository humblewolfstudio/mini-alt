package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/utils"
)

// GetObject gets the bucket and the file key and returns the file.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (h *Handler) GetObject(c *gin.Context, bucket string, objectKey string) {
	clientIp := utils.GetClientIP(c.Request)
	accessKey := c.GetString("accessKey")
	if accessKey == "" {
		utils.HandleError(c, utils.InternalServerError, "Access key not found")
		return
	}

	path, object, e := h.Storage.GetObject(bucket, objectKey)
	if e != "" {
		utils.HandleError(c, e, bucket)
		return
	}

	go events.HandleEventObject(h.Store, types.EventGet, utils.ClearInput(bucket), utils.ClearInput(objectKey), object.ETag, object.Size, accessKey, clientIp)

	c.File(path)
}
