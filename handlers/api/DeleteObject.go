package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/events"
	"mini-alt/events/types"
	"mini-alt/utils"
	"net/http"
)

// DeleteObject receives the key of the file and removes that file.
// If no file is found, it does not return an error, it just returns.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (h *Handler) DeleteObject(c *gin.Context, bucket, objectKey string) {
	ok, e := h.Storage.DeleteObject(bucket, objectKey)
	if !ok {
		utils.HandleError(c, e, bucket)
		return
	}

	go events.HandleEventObject(h.Store, types.EventDelete, utils.ClearObjectKeyWithBucket(bucket, objectKey), bucket, "")

	c.Status(http.StatusNoContent)
}
