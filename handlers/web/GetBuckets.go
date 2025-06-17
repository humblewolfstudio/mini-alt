package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
)

func createTestClient() *s3.S3 {
	cfg := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("http://localhost:8080"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			"YOUR_ACCESS_KEY_ID",
			"YOUR_SECRET_ACCESS_KEY",
			""),
	}

	sess := session.Must(session.NewSession(cfg))
	return s3.New(sess)
}

func ListBuckets(c *gin.Context) {
	s3Client := createTestClient()

	bucketList, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, bucketList.Buckets)
}
