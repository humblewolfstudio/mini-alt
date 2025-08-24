package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/utils"
	"net/http"
)

// HeadBucket returns if a buckets exists.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadBucket.html
func (h *Handler) HeadBucket(c *gin.Context) {
	clientIp := utils.GetClientIP(c.Request)
	accessKey := c.GetString("accessKey")
	if accessKey == "" {
		utils.HandleError(c, utils.InternalServerError, "Access key not found")
		return
	}

	bucket := utils.ClearInput(c.Param("bucket"))
	ok, e := h.Storage.HeadBucket(bucket)
	if !ok {
		utils.HandleError(c, e, bucket)
		return
	}

	go events.HandleEventBucket(h.Store, types.EventHead, bucket, accessKey, clientIp)

	c.Status(http.StatusOK)
}
