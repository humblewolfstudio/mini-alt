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
func (h *Handler) CreateBucket(c *gin.Context, bucketName string) {
	user, ok := GetUserFromContext(c)
	if !ok {
		utils.RespondS3Error(c, 500, "InternalServerError", "Error retrieving user", bucketName)
		return
	}

	if err := h.Store.PutBucket(bucketName, user.Id); err != nil {
		utils.RespondS3Error(c, http.StatusConflict, "BucketAlreadyExists",
			"The requested bucket name is not available.", bucketName)
		return
	}

	if err := disk.CreateBucket(bucketName); err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create storage directory.", bucketName)
		return
	}

	go events.HandleEventBucket(h.Store, types.EventBucketCreated, utils.ClearBucketName(bucketName), "")

	c.Status(http.StatusCreated)
}
