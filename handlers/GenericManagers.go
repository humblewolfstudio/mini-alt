package handlers

import "github.com/gin-gonic/gin"

// GetObjectOrList receives the endpoint of getting an object or listing bucket objects (due to gin problem with * endpoints).
func (h *Handler) GetObjectOrList(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")

	if object == "/" || object == "" {
		h.ListObjectsV2(c, bucket)
		return
	}

	object = object[1:]

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

	object = object[1:]
	h.PutObject(c, bucket, object)
}

func (h *Handler) DeleteObjectOrBucket(c *gin.Context) {
	object := c.Param("object")

	if object == "/" || object == "" {
		h.DeleteBucket(c)
	}

	object = object[1:]
	h.DeleteObject(c)
}
