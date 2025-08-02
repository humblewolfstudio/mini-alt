package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

func (h *Handler) UploadFiles(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	bucket := c.PostForm("bucket")
	prefix := c.PostForm("prefix")

	if bucket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket name is required"})
		return
	}

	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	s3Client := createTestClient(h, id.(int64))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		filename := filepath.Base(fileHeader.Filename)
		sanitizedFilename := strings.ReplaceAll(filename, " ", "_")

		key := path.Clean(prefix + sanitizedFilename)

		if strings.Contains(key, "..") || strings.HasPrefix(key, "/") || strings.Contains(key, "//") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object key path"})
			return
		}

		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   file,
		})

		_ = file.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "Upload failed",
				"detail": err.Error(),
				"key":    key,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Files uploaded successfully"})
}
