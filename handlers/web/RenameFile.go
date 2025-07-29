package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RenameFileRequest struct {
	Bucket string `json:"bucket" binding:"required"`
	OldKey string `json:"oldKey" binding:"required"`
	NewKey string `json:"newKey" binding:"required"`
}

func (h *WebHandler) RenameFile(c *gin.Context) {
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

	_, err := s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(req.Bucket),
		CopySource: aws.String(req.Bucket + "/" + req.OldKey),
		Key:        aws.String(req.NewKey),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(req.Bucket),
		Key:    aws.String(req.OldKey),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File renamed successfully"})
}
