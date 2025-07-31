package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/encoding"
	"mini-alt/storage/disk"
	"mini-alt/utils"
	"net/http"
	"strings"
)

// CopyObject
// TODO clean and add description
func (h *Handler) CopyObject(c *gin.Context, bucketName, objectKey, copySource string) {
	parts := strings.SplitN(copySource, "/", 2)
	/*
		if len(parts) != 2 {
		        http.Error(w, "Invalid x-amz-copy-source header", http.StatusBadRequest)
		        return
		    }
	*/
	srcBucketName := parts[0]
	srcObjectKey := parts[1]
	objectKey = strings.TrimPrefix(objectKey, "/")

	srcFile, err := disk.GetObject(srcBucketName, srcObjectKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "ObjectNotFound", "Could not find the Object.", srcBucketName)
		return
	}

	endPath, err := disk.GetObjectPath(bucketName, objectKey)

	written, err := disk.CreateObject(endPath, srcFile)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "CouldNotWrite", "Could not write the Object.", bucketName)
		return
	}

	object, err := h.Store.PutObject(bucketName, objectKey, written)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "CouldNotWrite", "Could not write the Object.", bucketName)
		return
	}

	var xmlCopyObjectResult encoding.CopyObjectResult

	xmlCopyObjectResult.LastModified = object.LastModified

	c.XML(http.StatusOK, xmlCopyObjectResult)
}
