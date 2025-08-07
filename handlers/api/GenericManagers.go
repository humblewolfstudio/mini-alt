package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/utils"
)

// GetObjectOrList receives the endpoint of getting an object or listing bucket objects (due to gin problem with * endpoints).
func (h *Handler) GetObjectOrList(c *gin.Context) {
	bucket := utils.ClearInput(c.Param("bucket"))
	objectKey := utils.ClearInput(c.Param("object"))

	if objectKey == "/" || objectKey == "" {
		h.ListObjectsV2(c, bucket)
		return
	}

	h.GetObject(c, bucket, objectKey)
}

// PutObjectOrBucket receives the endpoint of creating an object or a bucket (due to gin problem with * endpoints).
func (h *Handler) PutObjectOrBucket(c *gin.Context) {
	bucket := utils.ClearInput(c.Param("bucket"))
	objectKey := utils.ClearInput(c.Param("object"))

	if copySource := utils.ClearInput(c.GetHeader("x-amz-copy-source")); copySource != "" {
		h.CopyObject(c, bucket, objectKey, copySource)
		return
	}

	if objectKey == "/" || objectKey == "" {
		h.CreateBucket(c, bucket)
		return
	}

	h.PutObject(c, bucket, objectKey)
}

// DeleteObjectOrBucket receives the endpoint of deleting an object or a bucket (due to gin problem with * endpoints).
func (h *Handler) DeleteObjectOrBucket(c *gin.Context) {
	bucket := utils.ClearInput(c.Param("bucket"))
	objectKey := utils.ClearInput(c.Param("object"))

	if objectKey == "/" || objectKey == "" {
		h.DeleteBucket(c)
		return
	}

	h.DeleteObject(c, bucket, objectKey)
}

// HeadObjectOrBucket receives the endpoint of returning the metadata of an object or a bucket.
func (h *Handler) HeadObjectOrBucket(c *gin.Context) {
	bucket := utils.ClearInput(c.Param("bucket"))
	objectKey := utils.ClearInput(c.Param("object"))

	if objectKey == "/" || objectKey == "" {
		h.HeadBucket(c)
		return
	}

	h.HeadObject(c, bucket, objectKey)
}
