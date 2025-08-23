package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/utils"
	"net/http"
	"strconv"
)

// HeadObject returns the metadata of an object.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadObject.html
func (h *Handler) HeadObject(c *gin.Context, bucket, objectKey string) {
	object, metadata, e := h.Storage.HeadObject(bucket, objectKey)
	if e != "" {
		utils.HandleError(c, e, bucket)
		return
	}

	go events.HandleEventObject(h.Store, types.EventHead, utils.ClearObjectKeyWithBucket(bucket, objectKey), bucket, "")

	c.Header("Last-Modified", object.LastModified.Format(http.TimeFormat))
	c.Header("Content-Length", strconv.FormatInt(metadata.ContentLength, 10))
	c.Header("Content-Disposition", metadata.ContentDisposition)
	c.Header("Content-Encoding", metadata.ContentEncoding)
	c.Header("Content-Language", metadata.ContentLanguage)
	c.Header("Cache-Control", metadata.CacheControl)
	c.Header("Content-Type", metadata.ContentType)

	c.Status(http.StatusOK)
}
