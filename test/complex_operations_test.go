package test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
	"testing"
)

func renameObject(s3Client *s3.S3, bucket, srcKey, dstKey string) error {
	_, err := s3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(bucket + "/" + srcKey),
	})
	if err != nil {
		return err
	}

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
	})
	return err
}

func moveObject(s3Client *s3.S3, bucket, srcKey, dstKey string) error {
	return renameObject(s3Client, bucket, srcKey, dstKey)
}

func TestRenameObjectSuccess(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcKey := "file1.txt"
	dstKey := "file2.txt"

	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
		Body:   strings.NewReader("content"),
	})

	if err := renameObject(s3Client, bucket, srcKey, dstKey); err != nil {
		t.Fatalf("RenameObject failed: %v", err)
	}

	if _, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dstKey),
	}); err != nil {
		t.Errorf("Renamed object not found: %v", err)
	}

	if _, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
	}); err == nil {
		t.Errorf("Expected source to be deleted after rename")
	}
}

func TestRenameObjectNotFound(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	err := renameObject(s3Client, bucket, "does-not-exist.txt", "new.txt")
	if err == nil {
		t.Errorf("Expected error when renaming non-existent object")
	}
}

func TestMoveObjectOverwrite(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcKey := "file.txt"
	dstKey := "file-moved.txt"

	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
		Body:   strings.NewReader("original"),
	})
	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dstKey),
		Body:   strings.NewReader("existing"),
	})

	if err := moveObject(s3Client, bucket, srcKey, dstKey); err != nil {
		t.Fatalf("MoveObject failed: %v", err)
	}
}

func TestMoveObjectSameKeyFails(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcKey := "file.txt"
	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
		Body:   strings.NewReader("data"),
	})

	err := moveObject(s3Client, bucket, srcKey, srcKey)
	if err == nil {
		t.Errorf("Expected error when moving object to itself")
	}
}

func renameDirectory(s3Client *s3.S3, bucket, srcDir, dstDir string) error {
	list, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(srcDir),
	})
	if err != nil {
		return err
	}

	if len(list.Contents) == 0 {
		return fmt.Errorf("directory %s not found", srcDir)
	}

	for _, obj := range list.Contents {
		newKey := strings.Replace(*obj.Key, srcDir, dstDir, 1)
		_, err := s3Client.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(newKey),
			CopySource: aws.String(bucket + "/" + *obj.Key),
		})
		if err != nil {
			return err
		}
		_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    obj.Key,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func TestRenameDirectorySuccess(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcDir := "old/"
	dstDir := "new/"
	fileKey := srcDir + "file.txt"

	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		Body:   strings.NewReader("data"),
	})

	if err := renameDirectory(s3Client, bucket, srcDir, dstDir); err != nil {
		t.Fatalf("RenameDirectory failed: %v", err)
	}

	listOld, _ := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(srcDir),
	})
	if len(listOld.Contents) > 0 {
		t.Errorf("Expected old directory to be gone")
	}

	listNew, _ := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(dstDir),
	})
	if len(listNew.Contents) == 0 {
		t.Errorf("Expected new directory to contain objects")
	}
}

func TestRenameDirectoryNotFound(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	err := renameDirectory(s3Client, bucket, "missing/", "renamed/")
	if err == nil {
		t.Errorf("Expected error when renaming non-existent directory")
	}
}

func TestMoveDirectoryIntoItselfFails(t *testing.T) {
	s3Client := createTestClient()
	bucket := createTempBucket(t, s3Client)
	defer deleteTempBucket(t, s3Client, bucket)

	srcDir := "dir/"
	_, _ = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcDir + "file.txt"),
		Body:   strings.NewReader("data"),
	})

	err := renameDirectory(s3Client, bucket, srcDir, srcDir)
	if err == nil {
		t.Errorf("Expected error when moving directory into itself")
	}
}
