package test

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"mini-alt/utils"
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
			"0GNTGAiAxRywL5KI",
			"zj9VXm1q9bE8n9OEsgglpqec9DBtkkZe",
			""),
	}
	sess := session.Must(session.NewSession(cfg))
	return s3.New(sess)
}

func createTempBucket(t *testing.T, s3Client *s3.S3) string {
	bucket := fmt.Sprintf("test-bucket-%s", utils.GenerateRandomKey(16))
	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Fatalf("CreateBucket failed: %v", err)
	}
	return bucket
}

func deleteAllObjects(t *testing.T, s3Client *s3.S3, bucket string) {
	listOutput, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Fatalf("ListObjects failed: %v", err)
	}
	for _, obj := range listOutput.Contents {
		_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    obj.Key,
		})
		if err != nil {
			t.Errorf("DeleteObject %s failed: %v", *obj.Key, err)
		}
	}
}

func deleteTempBucket(t *testing.T, s3Client *s3.S3, bucket string) {
	deleteAllObjects(t, s3Client, bucket)
	_, err := s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Fatalf("DeleteBucket failed: %v", err)
	}
}

func TestPutAndGetObjectSuccess(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	key := "test-object.txt"
	content := "Hello, isolated S3 test!"
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
	})
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}

	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, result.Body)
	if err != nil {
		t.Fatalf("Read Body failed: %v", err)
	}
	if buf.String() != content {
		t.Errorf("Expected %q, got %q", content, buf.String())
	}
}

func TestGetObjectNotFound(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	_, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("does-not-exist"),
	})
	if err == nil {
		t.Errorf("Expected GetObject to fail for non-existent key")
	}
}

func TestPutObjectInvalidBucket(t *testing.T) {
	s3Client := createTestClient()
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("missing-bucket"),
		Key:    aws.String("x"),
		Body:   strings.NewReader("data"),
	})
	if err == nil {
		t.Errorf("Expected error when putting object in missing bucket")
	}
}

func TestCopyObjectSuccess(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcKey := "original.txt"
	dstKey := "copied.txt"

	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
		Body:   strings.NewReader("copy this"),
	})

	_, err := s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(bucket + "/" + srcKey),
	})
	if err != nil {
		t.Fatalf("CopyObject failed: %v", err)
	}
}

func TestCopyObjectSourceNotFound(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	_, err := s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String("dest.txt"),
		CopySource: aws.String(bucket + "/does-not-exist"),
	})
	if err == nil {
		t.Errorf("Expected error when copying missing object")
	}
}

func TestListObjectsWithPrefix(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("prefix/test1.txt"),
		Body:   strings.NewReader("data1"),
	})
	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("prefix/test2.txt"),
		Body:   strings.NewReader("data2"),
	})

	list, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String("prefix/"),
	})
	if err != nil {
		t.Fatalf("ListObjects failed: %v", err)
	}
	if len(list.Contents) != 2 {
		t.Errorf("Expected 2 objects, got %d", len(list.Contents))
	}
}

func TestListObjectsEmptyBucket(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	list, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Fatalf("ListObjects failed: %v", err)
	}
	if len(list.Contents) != 0 {
		t.Errorf("Expected 0 objects, got %d", len(list.Contents))
	}
}

func TestHeadObjectSuccess(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	key := "head-me.txt"
	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader("data"),
	})

	_, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("HeadObject failed: %v", err)
	}
}

func TestHeadObjectNotFound(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	_, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("not-there"),
	})
	if err == nil {
		t.Errorf("Expected HeadObject to fail for non-existent key")
	}
}

func TestDeleteObjectSuccess(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	key := "delete-me.txt"
	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader("temp"),
	})

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}

	_, err = s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		t.Errorf("Expected HeadObject to fail after deletion")
	}
}

func TestDeleteObjectNotFoundAllowed(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("does-not-exist"),
	})
	if err != nil {
		t.Errorf("Expected delete of missing object to succeed, got: %v", err)
	}
}
