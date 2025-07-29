package web

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type CreateFolderRequest struct {
	Bucket     string `json:"bucket"`
	Prefix     string `json:"prefix"`
	FolderName string `json:"folderName"`
}

func (h *WebHandler) CreateFolder(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	var req CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s3Client := createTestClient(h, id.(int64))

	key := strings.TrimSuffix(req.Prefix, "/") + "/" + strings.TrimPrefix(req.FolderName, "/") + "/"

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(req.Bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(""),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create folder",
			"detail":  err.Error(),
			"s3_path": fmt.Sprintf("s3://%s/%s", req.Bucket, key),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Folder created successfully",
		"path":    key,
	})
}
