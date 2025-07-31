package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage/disk"
	"mini-alt/utils"
	"net/http"
)

// CreateBucket receives the name of the new bucket and creates the bucket.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html
func (h *Handler) CreateBucket(c *gin.Context, bucketName string) {
	if err := h.Store.PutBucket(bucketName); err != nil {
		utils.RespondS3Error(c, http.StatusConflict, "BucketAlreadyExists",
			"The requested bucket name is not available.", bucketName)
		return
	}

	if err := disk.CreateBucketDirectory(bucketName); err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create storage directory.", bucketName)
		return
	}

	c.Status(http.StatusCreated)
}
