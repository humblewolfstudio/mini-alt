package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type MoveFileRequest struct {
	Bucket          string `json:"bucket" binding:"required"`
	SourceKey       string `json:"sourceKey" binding:"required"`
	DestinationPath string `json:"destinationPath" binding:"required"`
}

func (h *WebHandler) MoveFile(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	var req MoveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s3Client := createTestClient(h, id.(int64))

	filename := filepath.Base(req.SourceKey)
	newKey := req.DestinationPath + filename

	_, err := s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(req.Bucket),
		CopySource: aws.String(req.Bucket + "/" + req.SourceKey),
		Key:        aws.String(newKey),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(req.Bucket),
		Key:    aws.String(req.SourceKey),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File moved successfully"})
}
