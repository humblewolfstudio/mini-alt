package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/storage/disk"
	"mini-alt/utils"
	"net/http"
)

// DeleteBucket receives the bucket name and removes that bucket with all of its content
// If no bucket is found (or the bucket name is incorrect) it throws an error.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html
func (h *Handler) DeleteBucket(c *gin.Context) {
	bucket := utils.ClearInput(c.Param("bucket"))
	if err := disk.DeleteBucket(bucket); err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError", "Could not delete bucket", bucket)
		return
	}

	err := h.Store.DeleteBucket(bucket)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError", "Could not delete bucket", bucket)
		return
	}

	go events.HandleEventBucket(h.Store, types.EventBucketDeleted, bucket, "")

	c.Status(http.StatusNoContent)
}
