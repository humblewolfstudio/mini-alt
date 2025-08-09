package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/storage/disk"
	"net/http"
)

// CreateBucket receives the name of the new bucket and creates the bucket.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html
func (h *Handler) CreateBucket(c *gin.Context, bucket string) {
	if err := h.Store.PutBucket(bucket); err != nil {
		handleError(c, BucketAlreadyExists, bucket)
		return
	}

	if err := disk.CreateBucket(bucket); err != nil {
		handleError(c, FailedToCreateBucket, bucket)
		return
	}

	go events.HandleEventBucket(h.Store, types.EventBucketCreated, bucket, "")

	c.Status(http.StatusCreated)
}
