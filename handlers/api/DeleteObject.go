package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/storage/disk"
	"mini-alt/utils"
	"net/http"
)

// DeleteObject receives the key of the file and removes that file.
// If no file is found, it does not return an error, it just returns.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (h *Handler) DeleteObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	object := c.Param("object")

	// Also delete all files
	err := h.Store.DeleteObject(bucketName, object)
	if err == nil {
		err := disk.DeleteObjectFile(bucketName, object)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err != nil {
		println(err.Error())
	}

	go events.HandleEventObject(h.Store, types.EventDelete, utils.ClearObjectKeyWithBucket(bucketName, object), utils.ClearBucketName(bucketName), "")

	c.Status(http.StatusNoContent)
}
