package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createTestClient() *s3.S3 {
	cfg := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("http://localhost:9000"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			"0JIZzAUcBOfKIxjA",
			"G1EuwK4RKfL9cmUVuAMtvotxpTOXZDfn",
			""),
	}

	sess := session.Must(session.NewSession(cfg))
	return s3.New(sess)
}

type FileItem struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified"`
	IsFolder     bool   `json:"isFolder"`
}
