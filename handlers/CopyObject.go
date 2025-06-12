package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/encoding"
	"mini-alt/storage"
	"mini-alt/utils"
	"net/http"
	"strings"
)

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

	srcFile, err := storage.GetObject(srcBucketName, srcObjectKey)
	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "ObjectNotFound", "Could not find the Object.", srcBucketName)
	}

	endPath, err := storage.GetObjectPath(bucketName, objectKey)

	written, err := storage.CreateObject(endPath, srcFile)
	if err != nil {
		utils.RespondS3Error(c, http.StatusInternalServerError, "CouldNotWrite", "Could not write the Object.", bucketName)
	}

	object := h.Store.PutObject(bucketName, objectKey, written)

	var xmlCopyObjectResult encoding.CopyObjectResult

	xmlCopyObjectResult.LastModified = object.LastModified

	c.XML(http.StatusOK, xmlCopyObjectResult)
}
