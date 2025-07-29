package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *WebHandler) ListFolders(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	bucket := c.Query("bucket")
	excludePrefix := c.Query("excludePrefix")
	currentPath := c.Query("currentPath")

	s3Client := createTestClient(h, id.(int64))

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Delimiter: aws.String("/"),
	}

	if currentPath != "" {
		input.Prefix = aws.String(currentPath)
	}

	resp, err := s3Client.ListObjectsV2(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var folders []string

	if currentPath == "" {
		folders = append(folders, "")
	}

	for _, folder := range resp.CommonPrefixes {
		folderPath := *folder.Prefix

		if excludePrefix != "" && folderPath == excludePrefix {
			continue
		}
		folders = append(folders, folderPath)
	}

	c.JSON(http.StatusOK, gin.H{"data": folders})
}
