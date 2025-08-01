package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func (h *Handler) ListFiles(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	bucket := c.Query("bucket")
	prefix := c.Query("prefix")

	prefix = strings.TrimPrefix(prefix, "/")
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	s3Client := createTestClient(h, id.(int64))

	resp, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var files []FileItem
	for _, folder := range resp.CommonPrefixes {
		folderPrefix := *folder.Prefix
		if !strings.HasSuffix(folderPrefix, "/") {
			folderPrefix += "/"
		}

		folderName := strings.TrimPrefix(folderPrefix, prefix)
		folderName = strings.Trim(folderName, "/")
		if folderName == "" {
			folderName = strings.TrimSuffix(filepath.Base(folderPrefix), "/")
		}

		files = append(files, FileItem{
			Key:      *folder.Prefix,
			Name:     folderName,
			IsFolder: true,
		})
	}

	for _, obj := range resp.Contents {
		if strings.HasSuffix(*obj.Key, "/") {
			continue
		}

		files = append(files, FileItem{
			Key:          *obj.Key,
			Name:         filepath.Base(*obj.Key),
			Size:         *obj.Size,
			LastModified: obj.LastModified.Format(time.RFC3339),
			IsFolder:     false,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": files})
}
