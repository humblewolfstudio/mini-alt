package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
	"testing"
)

func TestCreateBucketAlreadyExists(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err == nil {
		t.Errorf("expected error when creating existing bucket, got nil")
	}
}

func TestCreateBucketInvalidName(t *testing.T) {
	s3Client := createTestClient()

	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("INVALID_BUCKET_NAME!!!"),
	})
	if err == nil {
		t.Errorf("expected error for invalid bucket name, got nil")
	}
}

func TestHeadBucketExists(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	_, err := s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Errorf("HeadBucket failed on existing bucket: %v", err)
	}
}

func TestHeadBucketNotExists(t *testing.T) {
	s3Client := createTestClient()

	_, err := s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String("non-existent-bucket-12345"),
	})
	if err == nil {
		t.Errorf("expected error for HeadBucket on non-existent bucket, got nil")
	}
}

func TestDeleteBucketSuccess(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)

	_, err := s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Errorf("DeleteBucket failed: %v", err)
	}
}

func TestDeleteBucketNotExists(t *testing.T) {
	s3Client := createTestClient()

	_, err := s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String("non-existent-bucket-12345"),
	})
	if err == nil {
		t.Errorf("expected error deleting non-existent bucket, got nil")
	}
}

func TestDeleteBucketNotEmpty(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("file.txt"),
		Body:   strings.NewReader("data"),
	})

	_, err := s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err == nil {
		t.Errorf("expected error deleting non-empty bucket, got nil")
	}
}

func TestListBuckets(t *testing.T) {
	s3Client := createTestClient()
	_, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		t.Fatalf("ListBuckets failed: %v", err)
	}
}
