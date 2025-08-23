package storage

import (
	"io"
	"mini-alt/models"
	"mini-alt/storage/disk"
	"mini-alt/types"
	"mini-alt/utils"
)

func (s *Storage) putObject(bucket, objectKey string, body io.Reader, metadata types.Metadata, owner int64) (string, utils.Error) {
	_, err := s.store.GetBucket(bucket)
	if err != nil {
		return "", utils.NoSuchBucket
	}

	written, err := disk.PutObject(bucket, objectKey, body)
	if err != nil {
		return "", utils.InvalidRequest
	}

	object, err := s.store.PutObject(bucket, objectKey, written)
	if err != nil {
		return "", utils.InvalidRequest
	}

	err = s.store.PutMetadata(object.Id, metadata)
	if err != nil {
		return "", utils.InvalidRequest
	}

	md5, _ := disk.GetMD5Base64(bucket, objectKey)

	return md5, ""
}

func (s *Storage) deleteObject(bucket, objectKey string) (bool, utils.Error) {
	err := s.store.DeleteObject(bucket, objectKey)
	if err != nil {
		return false, utils.FailedToDeleteObject
	}

	err = disk.DeleteObject(bucket, objectKey)
	if err != nil {
		return false, utils.FailedToDeleteObjectFile
	}

	return true, ""
}

func (s *Storage) headObject(bucket, objectKey string) (*models.Object, *models.ObjectMetadata, utils.Error) {
	object, err := s.store.GetObject(bucket, objectKey)
	if err != nil {
		return nil, nil, utils.NoSuchKey
	}

	metadata, err := s.store.GetMetadata(object.Id)
	if err != nil {
		return nil, nil, utils.NoSuchKey
	}

	return &object, &metadata, ""
}

func (s *Storage) copyObject(srcBucket, srcKey, dstBucket, dstKey string) (*models.Object, utils.Error) {
	srcFile, err := disk.GetObject(srcBucket, srcKey)
	if err != nil {
		return nil, utils.NoSuchKey
	}

	written, err := disk.PutObject(dstBucket, dstKey, srcFile)
	if err != nil {
		return nil, utils.InvalidRequest
	}

	oldObject, err := s.store.GetObject(srcBucket, srcKey)
	if err != nil {
		return nil, utils.NoSuchKey
	}

	object, err := s.store.PutObject(dstBucket, dstKey, written)
	if err != nil {
		return nil, utils.InvalidRequest
	}

	err = s.store.CopyMetadata(oldObject.Id, object.Id)
	if err != nil {
		return nil, utils.InvalidRequest
	}

	return &object, ""
}

func (s *Storage) getObjectPath(bucket, objectKey string) (string, utils.Error) {
	path, err := disk.GetSafeObjectPath(bucket, objectKey)
	if err != nil {
		return "", utils.NoSuchKey
	}

	return path, ""
}
