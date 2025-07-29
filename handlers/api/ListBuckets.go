package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/encoding"
	"net/http"
	"time"
)

// ListBuckets returns a list of all buckets stored in memory.
// Each bucket contains its name and creation timestamp.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html
func (h *Handler) ListBuckets(c *gin.Context) {
	buckets, _ := h.Store.ListBuckets()
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
