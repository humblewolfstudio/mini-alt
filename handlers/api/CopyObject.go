package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/encoding"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/storage/disk"
	"mini-alt/utils"
	"net/http"
	"net/url"
	"strings"
)

// CopyObject
// TODO clean and add description
func (h *Handler) CopyObject(c *gin.Context, bucketName, objectKey, copySource string) {
	decodedCopySource, err := url.PathUnescape(copySource)
	if err != nil {
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidSourceKey", "Invalid source object key encoding", "")
	}

	parts := strings.SplitN(decodedCopySource, "/", 3)

	srcBucketName := parts[0]
	srcObjectKey := parts[1]

	if parts[0] == "" {
		srcBucketName = parts[1]
		srcObjectKey = parts[2]
	}

	decodedSrcKey, err := url.PathUnescape(srcObjectKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidSourceKey", "Invalid source object key encoding", srcBucketName)
		return
	}

	decodedDstKey, err := url.PathUnescape(strings.TrimPrefix(objectKey, "/"))
	if err != nil {
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidDestinationKey", "Invalid destination object key encoding", bucketName)
		return
	}

	srcFile, err := disk.GetObject(srcBucketName, decodedSrcKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "ObjectNotFound", "Could not find the Object.", srcBucketName)
		return
	}

	written, err := disk.PutObject(bucketName, decodedDstKey, srcFile)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "CouldNotWrite", "Could not write the Object.", bucketName)
		return
	}

	oldObject, err := h.Store.GetObject(srcBucketName, decodedSrcKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "ObjectNotFound", "Could not find the source Object.", bucketName)
		return
	}

	object, err := h.Store.PutObject(bucketName, decodedDstKey, written)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "CouldNotWrite", "Could not write the destination Object.", bucketName)
		return
	}

	err = h.Store.CopyMetadata(oldObject.Id, object.Id)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "CouldNotWrite", "Could not copy the metadata.", bucketName)
		return
	}

	var xmlCopyObjectResult encoding.CopyObjectResult

	xmlCopyObjectResult.LastModified = object.LastModified

	go events.HandleEventObject(h.Store, types.EventCopy, utils.ClearObjectKeyWithBucket(srcBucketName, decodedSrcKey), utils.ClearBucketName(srcBucketName), "")
	go events.HandleEventObject(h.Store, types.EventCopied, utils.ClearObjectKeyWithBucket(bucketName, decodedDstKey), utils.ClearBucketName(bucketName), "")

	c.XML(http.StatusOK, xmlCopyObjectResult)
}
