package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mini-alt/encoding"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/utils"
	"net/http"
	"net/url"
	"strings"
)

// CopyObject copies the object given a source bucket and key and a destination bucket and key
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CopyObject.html
func (h *Handler) CopyObject(c *gin.Context, dstBucket, objectKey, copySource string) {
	clientIp := utils.GetClientIP(c.Request)
	accessKey := c.GetString("accessKey")
	if accessKey == "" {
		utils.HandleError(c, utils.InternalServerError, "Access key not found")
		return
	}

	srcBucket, srcObjectKey, err := parseCopySource(copySource)
	if err != nil {
		utils.HandleError(c, utils.InvalidSourceKey, srcBucket)
		return
	}

	dstObjectKey, err := decodeObjectKey(objectKey)
	if err != nil {
		utils.HandleError(c, utils.InvalidDestinationKey, dstBucket)
		return
	}

	if srcObjectKey == dstObjectKey && srcBucket == dstBucket {
		utils.HandleError(c, utils.InvalidDestinationKey, srcBucket)
		return
	}

	object, errMsg := h.Storage.CopyObject(srcBucket, srcObjectKey, dstBucket, dstObjectKey)
	if errMsg != "" {
		utils.HandleError(c, errMsg, dstBucket)
		return
	}

	go events.HandleEventObjectCopy(h.Store, types.EventCopy, utils.ClearInput(srcBucket), utils.ClearInput(srcObjectKey), utils.ClearInput(dstBucket), utils.ClearInput(dstObjectKey), object.ETag, object.Size, accessKey, clientIp)

	c.XML(http.StatusOK, encoding.CopyObjectResult{
		LastModified: object.LastModified,
	})
}

func parseCopySource(copySource string) (string, string, error) {
	decoded, err := url.PathUnescape(copySource)
	if err != nil {
		return "", "", fmt.Errorf("invalid source object key encoding")
	}

	parts := strings.SplitN(decoded, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid copy source format")
	}

	srcBucket := parts[0]
	srcObjectKey := utils.ClearInput(parts[1])

	decodedSrcKey, err := url.PathUnescape(srcObjectKey)
	if err != nil {
		return "", "", fmt.Errorf("invalid source object key encoding")
	}

	return srcBucket, decodedSrcKey, nil
}

func decodeObjectKey(objectKey string) (string, error) {
	cleanedKey := utils.ClearInput(strings.TrimPrefix(objectKey, "/"))
	decoded, err := url.PathUnescape(cleanedKey)
	if err != nil {
		return "", fmt.Errorf("invalid object key encoding")
	}
	return decoded, nil
}
