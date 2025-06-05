package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/utils"
	"net/http"
	"os"
)

// CreateBucket receives the name of the new bucket and creates the bucket.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html
func (h *Handler) CreateBucket(c *gin.Context, bucket string) {
	if err := h.Store.CreateBucket(bucket); err != nil {
		utils.RespondS3Error(c, http.StatusConflict, "BucketAlreadyExists",
			"The requested bucket name is not available.", bucket)
		return
	}

	if err := os.MkdirAll("uploads/"+bucket, os.ModePerm); err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create storage directory.", bucket)
		return
	}

	c.Status(http.StatusCreated)
}
