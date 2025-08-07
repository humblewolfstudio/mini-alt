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
	bucketName := c.Param("bucket")

	_, err := h.Store.GetBucket(bucketName)
	if err != nil {
		println(err.Error())
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchBucket", err.Error(), bucketName)
		return
	}

	go events.HandleEventBucket(h.Store, types.EventHead, utils.ClearBucketName(bucketName), "")

	c.Status(http.StatusOK)
}
