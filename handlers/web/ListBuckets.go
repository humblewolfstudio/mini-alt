package web

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListBuckets(c *gin.Context) {
	s3Client := createTestClient()

	bucketList, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bucketList.Buckets)
}
