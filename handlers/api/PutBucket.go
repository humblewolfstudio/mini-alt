package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/utils"
	"net/http"
)

// PutBucket receives the name of the new bucket and creates the bucket.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html
func (h *Handler) PutBucket(c *gin.Context, bucket string) {
	clientIp := utils.GetClientIP(c.Request)
	accessKey := c.GetString("accessKey")
	if accessKey == "" {
		utils.HandleError(c, utils.InternalServerError, "Access key not found")
		return
	}

	user, ok := GetUserFromContext(c)
	if !ok {
		utils.HandleError(c, utils.InternalServerError, bucket)
		return
	}

	ok, e := h.Storage.PutBucket(bucket, user.Id)
	if !ok {
		utils.HandleError(c, e, bucket)
		return
	}

	go events.HandleEventBucket(h.Store, types.EventBucketCreated, bucket, accessKey, clientIp)

	c.Status(http.StatusCreated)
}
