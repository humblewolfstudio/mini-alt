package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/utils"
	"net/http"
	"strconv"
)

// HeadObject returns the metadata of an object.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadObject.html
func (h *Handler) HeadObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectKey := c.Param("object")

	object, err := h.Store.GetObject(bucketName, objectKey)

	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchKey", err.Error(), bucketName)
		return
	}

	c.Header("Last-Modified", object.LastModified.Format(http.TimeFormat))
	c.Header("Content-Length", strconv.FormatInt(object.Size, 10))

	c.Status(http.StatusOK)
}
