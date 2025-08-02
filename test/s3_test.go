package storage_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"strings"
	"testing"
)

func createTestClient() *s3.S3 {
	cfg := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("http://localhost:9000"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			"4yQlUGKo63I38E5S",
			"Cj50bBelkOI61lhC5clYcYAwqS5RzXPY",
			""),
	}

	sess := session.Must(session.NewSession(cfg))
	return s3.New(sess)
}

func TestPutObject(t *testing.T) {
	s3Client := createTestClient()

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-object.txt"),
		Body:   strings.NewReader("Hello, S3-compatible storage!"),
	})
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}
}

func TestGetObject(t *testing.T) {
	s3Client := createTestClient()

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-object.txt"),
		Body:   strings.NewReader("Hello, S3-compatible storage!"),
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-object.txt"),
	})
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(result.Body)

	body := make([]byte, 100)
	n, err := result.Body.Read(body)
	if err != nil && err != io.EOF {
		t.Fatalf("Read failed: %v", err)
	}

	expected := "Hello, S3-compatible storage!"
	if string(body[:n]) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(body[:n]))
	}
}

func TestCopyObject(t *testing.T) {
	s3Client := createTestClient()

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-object.txt"),
		Body:   strings.NewReader("Hello, S3-compatible storage!"),
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	_, err = s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String("test-bucket"),
		Key:        aws.String("test-dir/test-object-copy.txt"),
		CopySource: aws.String("test-bucket/test-object.txt"),
	})
	if err != nil {
		t.Fatalf("CopyObject failed: %v", err)
	}

	_, err = s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-dir/test-object-copy.txt"),
	})
	if err != nil {
		t.Errorf("Verification failed - copied object not found: %v", err)
	}
}

func TestListObjects(t *testing.T) {
	s3Client := createTestClient()

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-object.txt"),
		Body:   strings.NewReader("Hello, S3-compatible storage!"),
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("testDir/test-object.txt"),
		Body:   strings.NewReader("Hello, S3-compatible storage!"),
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	result, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String("test-bucket"),
		Prefix:    aws.String("testDir/"),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		t.Fatalf("ListObjects failed: %v", err)
	}

	found := false
	for _, item := range result.Contents {
		if *item.Key == "testDir/test-object.txt" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected object not found in listing")
	}
}

func TestHeadBucket(t *testing.T) {
	s3Client := createTestClient()

	_, err := s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String("test-bucket"),
	})

	if err != nil {
		t.Fatalf("HeadBucket failed: %v", err)
	}
}

func TestHeadObject(t *testing.T) {
	s3Client := createTestClient()

	_, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-object.txt"),
	})

	if err != nil {
		t.Fatalf("HeadBucket failed: %v", err)
	}
}

func TestDeleteBucket(t *testing.T) {
	s3Client := createTestClient()

	bucketName := "test-bucket"

	listOutput, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("Failed to list objects before deleting bucket: %v", err)
	}

	for _, obj := range listOutput.Contents {
		_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    obj.Key,
		})
		if err != nil {
			t.Errorf("Failed to delete object %s: %v", *obj.Key, err)
		}
	}

	err = s3Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("test-object.txt"),
	})
	if err != nil {
		t.Logf("WaitUntilObjectNotExists may have failed (possibly due to missing object): %v", err)
	}

	_, err = s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("DeleteBucket failed: %v", err)
	}
}
