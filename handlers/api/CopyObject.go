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

// CopyObject
// TODO clean and add description
func (h *Handler) CopyObject(c *gin.Context, dstBucket, objectKey, copySource string) {
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

	h.handleCopyEvents(srcBucket, srcObjectKey, dstBucket, dstObjectKey)

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

func (h *Handler) handleCopyEvents(srcBucket, srcKey, dstBucket, dstKey string) {
	srcEventKey := utils.ClearObjectKeyWithBucket(srcBucket, srcKey)
	dstEventKey := utils.ClearObjectKeyWithBucket(dstBucket, dstKey)

	go events.HandleEventObject(h.Store, types.EventCopy, srcEventKey, srcBucket, "")
	go events.HandleEventObject(h.Store, types.EventCopied, dstEventKey, dstBucket, "")
}
