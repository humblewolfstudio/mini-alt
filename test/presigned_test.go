package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

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
	_ = resp.Body.Close()

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
	_ = resp.Body.Close()

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
	bodyBytes, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

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
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed with status %d", resp.StatusCode)
	}

	headURL, err = generatePresignedHeadURL(bucket, key)
	if err != nil {
		t.Fatalf("Failed to generate HEAD URL: %v", err)
	}

	resp, err = http.Head(headURL)
	if err != nil {
		t.Fatalf("HEAD after delete request failed: %v", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected HEAD to fail after deletion, but got %d", resp.StatusCode)
	}
}

func TestPresignedUrlExpired(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	key := "expired.txt"
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(1 * time.Second)
	if err != nil {
		t.Fatalf("failed to presign request: %v", err)
	}

	time.Sleep(2 * time.Second)

	httpReq, _ := http.NewRequest("PUT", url, strings.NewReader("expired data"))
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Errorf("expected expired URL to fail, got status %d", resp.StatusCode)
	}
}

func TestPresignedWrongMethod(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	key := "wrong-method.txt"
	putURL, _ := generatePresignedPutURL(bucket, key)

	// Try GET with PUT URL
	resp, err := http.Get(putURL)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Errorf("expected failure when using GET on PUT presigned URL, got %d", resp.StatusCode)
	}
}

func TestPresignedGetNonExistentObject(t *testing.T) {
	bucket := "nonexistent-bucket"
	key := "does-not-exist.txt"

	url, _ := generatePresignedGetURL(bucket, key)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Errorf("expected error for non-existent object, got %d", resp.StatusCode)
	}
}

func TestPresignedDeleteNonExistentObject(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	key := "no-file.txt"
	url, _ := generatePresignedDeleteURL(bucket, key)
	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Errorf("expected error deleting non-existent object, got 200")
	}
}
