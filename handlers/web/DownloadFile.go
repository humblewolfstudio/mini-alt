package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"path/filepath"
)

func (h *WebHandler) DownloadFile(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	bucket := c.Query("bucket")
	key := c.Query("key")

	s3Client := createTestClient(h, id.(int64))

	resp, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Disposition", "attachment; filename=\""+filepath.Base(key)+"\"")
	if resp.ContentType != nil {
		c.Header("Content-Type", *resp.ContentType)
	} else {
		c.Header("Content-Type", "application/octet-stream")
	}

	c.Status(http.StatusOK)

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
}
