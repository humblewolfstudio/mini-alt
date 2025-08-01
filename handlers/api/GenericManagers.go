package api

import (
	"github.com/gin-gonic/gin"
)

// GetObjectOrList receives the endpoint of getting an object or listing bucket objects (due to gin problem with * endpoints).
func (h *Handler) GetObjectOrList(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")

	if object == "/" || object == "" {
		h.ListObjectsV2(c, bucket)
		return
	}

	h.GetObject(c, bucket, object)
}

// PutObjectOrBucket receives the endpoint of creating an object or a bucket (due to gin problem with * endpoints).
func (h *Handler) PutObjectOrBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")

	if copySource := c.GetHeader("x-amz-copy-source"); copySource != "" {
		h.CopyObject(c, bucket, object, copySource)
		return
	}

	if object == "/" || object == "" {
		h.CreateBucket(c, bucket)
		return
	}

	h.PutObject(c, bucket, object)
}

// DeleteObjectOrBucket receives the endpoint of deleting an object or a bucket (due to gin problem with * endpoints).
func (h *Handler) DeleteObjectOrBucket(c *gin.Context) {
	object := c.Param("object")

	if object == "/" || object == "" {
		h.DeleteBucket(c)
		return
	}

	h.DeleteObject(c)
}

// HeadObjectOrBucket receives the endpoint of returning the metadata of an object or a bucket.
func (h *Handler) HeadObjectOrBucket(c *gin.Context) {
	object := c.Param("object")

	if object == "/" || object == "" {
		h.HeadBucket(c)
		return
	}

	h.HeadObject(c)
}
