package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"sort"
)

type RenameFileRequest struct {
	Bucket string `json:"bucket" binding:"required"`
	OldKey string `json:"oldKey" binding:"required"`
	NewKey string `json:"newKey" binding:"required"`
}

func (h *Handler) RenameFile(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	var req RenameFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s3Client := createTestClient(h, id.(int64))
	var keys []*string

	err := s3Client.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(req.Bucket),
		Prefix: aws.String(req.OldKey),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			keys = append(keys, obj.Key)
		}
		return true
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "listing failed: " + err.Error()})
		return
	}

	sort.Slice(keys, func(i, j int) bool {
		return len(*keys[i]) > len(*keys[j])
	})

	for _, key := range keys {
		newObjKey := req.NewKey + (*key)[len(req.OldKey):]
		copySource := url.PathEscape("/" + req.Bucket + "/" + *key)

		_, err := s3Client.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(req.Bucket),
			CopySource: aws.String(copySource),
			Key:        aws.String(newObjKey),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "copy failed: " + err.Error()})
			return
		}

		_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(req.Bucket),
			Key:    key,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Renamed successfully"})
}
