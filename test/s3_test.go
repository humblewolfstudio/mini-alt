package storage_test

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"mini-alt/utils"
	"net/http"
	"strings"
	"testing"
	"time"
)

func createTestClient() *s3.S3 {
	//goland:noinspection SpellCheckingInspection
	cfg := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("http://localhost:9000"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			"JXFLFwjme1d31Fe8",
			"5JQniImyOxsoadwQxju3SkqQ6DdhQbxg",
			"",
		),
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

func TestPutAndGetObject(t *testing.T) {
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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(result.Body)

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, result.Body)
	if err != nil {
		t.Fatalf("Read Body failed: %v", err)
	}
	if buf.String() != content {
		t.Errorf("Expected %q, got %q", content, buf.String())
	}
}

func TestCopyObject(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcKey := "original.txt"
	dstKey := "copied.txt"

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
		Body:   strings.NewReader("copy this"),
	})
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}

	_, err = s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(bucket + "/" + srcKey),
	})
	if err != nil {
		t.Fatalf("CopyObject failed: %v", err)
	}

	_, err = s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dstKey),
	})
	if err != nil {
		t.Fatalf("Copied object not found: %v", err)
	}
}

func TestListObjects(t *testing.T) {
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

func TestHeadBucket(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	_, err := s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Fatalf("HeadBucket failed: %v", err)
	}
}

func TestHeadObject(t *testing.T) {
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

func TestDeleteObject(t *testing.T) {
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

func TestListBuckets(t *testing.T) {
	s3Client := createTestClient()
	_, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		t.Fatalf("ListBuckets failed: %v", err)
	}
}

func TestRenameFile(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcKey := "original-file.txt"
	dstKey := "renamed-file.txt"

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
		Body:   strings.NewReader("file content"),
	})
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}

	_, err = s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(bucket + "/" + srcKey),
	})
	if err != nil {
		t.Fatalf("CopyObject failed: %v", err)
	}

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
	})
	if err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}

	_, err = s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dstKey),
	})
	if err != nil {
		t.Errorf("Renamed object not found: %v", err)
	}
}

func TestMoveFileToDirectory(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcKey := "file.txt"
	dstKey := "subdir/file.txt"

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
		Body:   strings.NewReader("hello move"),
	})
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}

	_, err = s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(bucket + "/" + srcKey),
	})
	if err != nil {
		t.Fatalf("CopyObject failed: %v", err)
	}

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
	})
	if err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}

	_, err = s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dstKey),
	})
	if err != nil {
		t.Errorf("Moved object not found: %v", err)
	}
}

func TestRenameDirectory(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcDir := "old-dir/"
	dstDir := "new-dir/"

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcDir),
		Body:   strings.NewReader(""),
	})
	if err != nil {
		t.Fatalf("PutObject (directory) failed: %v", err)
	}

	_, err = s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(dstDir),
		CopySource: aws.String(bucket + "/" + srcDir),
	})
	if err != nil {
		t.Fatalf("CopyObject (directory) failed: %v", err)
	}

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcDir),
	})
	if err != nil {
		t.Fatalf("DeleteObject (directory) failed: %v", err)
	}

	_, err = s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dstDir),
	})
	if err != nil {
		t.Errorf("Renamed directory not found: %v", err)
	}
}

func TestMoveDirectory(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcDir := "folder1/"
	dstDir := "folder2/"
	fileName := "file.txt"
	srcKey := srcDir + fileName
	dstKey := dstDir + fileName

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
		Body:   strings.NewReader("folder content"),
	})
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}

	_, err = s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(bucket + "/" + srcKey),
	})
	if err != nil {
		t.Fatalf("CopyObject failed: %v", err)
	}

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
	})
	if err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}

	_, err = s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dstKey),
	})
	if err != nil {
		t.Errorf("Moved file not found in destination directory: %v", err)
	}
}

func generatePresignedPutURL(bucket, key string) (string, error) {
	s3Client := createTestClient()
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(5 * time.Minute)
}

func generatePresignedGetURL(bucket, key string) (string, error) {
	s3Client := createTestClient()
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(5 * time.Minute)
}

func generatePresignedHeadURL(bucket, key string) (string, error) {
	s3Client := createTestClient()
	req, _ := s3Client.HeadObjectRequest(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(5 * time.Minute)
}

func generatePresignedDeleteURL(bucket, key string) (string, error) {
	s3Client := createTestClient()
	req, _ := s3Client.DeleteObjectRequest(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(5 * time.Minute)
}

func TestPresignedObjectOperations(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	key := "test-presigned.txt"
	content := "presigned test content"

	putURL, err := generatePresignedPutURL(bucket, key)
	if err != nil {
		t.Fatalf("Failed to generate PUT URL: %v", err)
	}

	req, err := http.NewRequest("PUT", putURL, strings.NewReader(content))
	if err != nil {
		t.Fatalf("Failed to create PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT failed with status %d", resp.StatusCode)
	}

	headURL, err := generatePresignedHeadURL(bucket, key)
	if err != nil {
		t.Fatalf("Failed to generate HEAD URL: %v", err)
	}

	resp, err = http.Head(headURL)
	if err != nil {
		t.Fatalf("HEAD request failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HEAD failed with status %d", resp.StatusCode)
	}

	getURL, err := generatePresignedGetURL(bucket, key)
	if err != nil {
		t.Fatalf("Failed to generate GET URL: %v", err)
	}

	resp, err = http.Get(getURL)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read GET response: %v", err)
	}

	if string(bodyBytes) != content {
		t.Errorf("GET content mismatch: expected %q, got %q", content, string(bodyBytes))
	}

	deleteURL, err := generatePresignedDeleteURL(bucket, key)
	if err != nil {
		t.Fatalf("Failed to generate DELETE URL: %v", err)
	}

	req, err = http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed with status %d", resp.StatusCode)
	}

	_, err = s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		t.Errorf("Expected object to be deleted, but HEAD still succeeds")
	}
}
