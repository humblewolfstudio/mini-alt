package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage"
	"net/http"
)

// DeleteObject receives the key of the file and removes that file.
// If no file is found, it does not return an error, it just returns.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (h *Handler) DeleteObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectKey := c.Param("object")
	object := objectKey[1:]

	// Also delete all files
	err := h.Store.DeleteObject(bucketName, object)
	if err == nil {
		err := storage.DeleteObjectFile(bucketName, object)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err != nil {
		println(err.Error())
	}

	println(err)
	c.Status(http.StatusNoContent)
}
