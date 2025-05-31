package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/encoding"
	"mini-alt/utils"
	"net/http"
	"os"
	"time"
)

func (h *Handler) ListBuckets(c *gin.Context) {
	buckets := h.Store.ListBuckets()
	var xmlBuckets []encoding.BucketXML

	for _, bucket := range buckets {
		xmlBuckets = append(xmlBuckets, encoding.BucketXML{
			Name:      bucket.Name,
			CreatedAt: bucket.CreatedAt.Format(time.RFC3339),
		})
	}

	res := encoding.ListAllMyBucketsResult{}
	res.Buckets.Bucket = xmlBuckets
	c.XML(http.StatusOK, res)
}

func (h *Handler) handleCreateBucket(c *gin.Context, bucket string) {
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
