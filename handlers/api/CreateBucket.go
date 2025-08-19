package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/storage/disk"
	"mini-alt/utils"
	"net/http"
)

// CreateBucket receives the name of the new bucket and creates the bucket.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html
func (h *Handler) CreateBucket(c *gin.Context, bucket string) {
	user, ok := GetUserFromContext(c)
	if !ok {
		utils.RespondS3Error(c, 500, "InternalServerError", "Error retrieving user", bucket)
		return
	}

	if err := h.Store.PutBucket(bucket, user.Id); err != nil {
		handleError(c, BucketAlreadyExists, bucket)
	}

	if err := disk.CreateBucket(bucket); err != nil {
		handleError(c, FailedToCreateBucket, bucket)
		return
	}

	go events.HandleEventBucket(h.Store, types.EventBucketCreated, bucket, "")

	c.Status(http.StatusCreated)
}
