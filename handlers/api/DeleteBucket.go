package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/utils"
	"net/http"
)

// DeleteBucket receives the bucket name and removes that bucket with all of its content
// If no bucket is found (or the bucket name is incorrect) it throws an error.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html
func (h *Handler) DeleteBucket(c *gin.Context) {
	clientIp := utils.GetClientIP(c.Request)
	accessKey := c.GetString("accessKey")
	if accessKey == "" {
		utils.HandleError(c, utils.InternalServerError, "Access key not found")
		return
	}

	bucket := utils.ClearInput(c.Param("bucket"))
	ok, e := h.Storage.DeleteBucket(bucket)
	if !ok {
		utils.HandleError(c, e, bucket)
		return
	}

	go events.HandleEventBucket(h.Store, types.EventBucketDeleted, bucket, accessKey, clientIp)

	c.Status(http.StatusNoContent)
}
