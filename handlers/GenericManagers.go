package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) GetObjectCommandOrList(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")

	if object == "/" || object == "" {
		h.ListObjectsV2(c, bucket)
		return
	}

	object = object[1:]

	h.GetObject(c, bucket, object)
}
