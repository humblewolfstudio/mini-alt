package web

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mini-alt/storage"
)

func createTestClient(h *WebHandler, id int64) *s3.S3 {
	creds, err := getCredentials(h, id)
	if err != nil {
		return nil
	}

	cfg := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("http://localhost:9000"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			creds.AccessKey,
			creds.SecretKey,
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

func getCredentials(h *WebHandler, id int64) (storage.Credentials, error) {
	user, err := h.Store.GetUserById(id)
	if err != nil {
		return storage.Credentials{}, err
	}

	secretKey, err := h.Store.GetSecretKey(user.AccessKey)
	if err != nil {
		return storage.Credentials{}, err
	}

	return storage.Credentials{AccessKey: user.AccessKey, SecretKey: secretKey}, nil
}
