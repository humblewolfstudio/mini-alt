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
func (h *Handler) CopyObject(c *gin.Context, bucket, objectKey, copySource string) {
	decodedCopySource, err := url.PathUnescape(copySource)
	if err != nil {
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidSourceKey", "Invalid source object key encoding", "")
	}

	parts := strings.SplitN(decodedCopySource, "/", 2)
	srcBucket := parts[0]
	srcObjectKey := utils.ClearInput(parts[1])

	decodedSrcKey, err := url.PathUnescape(srcObjectKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidSourceKey", "Invalid source object key encoding", srcBucket)
		return
	}

	decodedDstKey, err := url.PathUnescape(strings.TrimPrefix(objectKey, "/"))
	if err != nil {
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidDestinationKey", "Invalid destination object key encoding", bucket)
		return
	}

	srcFile, err := disk.GetObject(srcBucket, decodedSrcKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "ObjectNotFound", "Could not find the Object.", srcBucket)
		return
	}

	endPath, err := disk.GetObjectPath(bucket, decodedDstKey)

	written, err := disk.CreateObject(endPath, srcFile)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "CouldNotWrite", "Could not write the Object.", bucket)
		return
	}

	oldObject, err := h.Store.GetObject(srcBucket, decodedSrcKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "ObjectNotFound", "Could not find the source Object.", bucket)
		return
	}

	object, err := h.Store.PutObject(bucket, decodedDstKey, written)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "CouldNotWrite", "Could not write the destination Object.", bucket)
		return
	}

	err = h.Store.MetadataCopy(oldObject.Id, object.Id)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "CouldNotWrite", "Could not copy the metadata.", bucket)
		return
	}

	var xmlCopyObjectResult encoding.CopyObjectResult

	xmlCopyObjectResult.LastModified = object.LastModified

	go events.HandleEventObject(h.Store, types.EventCopy, utils.ClearObjectKeyWithBucket(srcBucket, decodedSrcKey), srcBucket, "")
	go events.HandleEventObject(h.Store, types.EventCopied, utils.ClearObjectKeyWithBucket(bucket, decodedDstKey), bucket, "")

	c.XML(http.StatusOK, xmlCopyObjectResult)
}
