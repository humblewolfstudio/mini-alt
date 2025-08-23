package utils

import (
	"github.com/gin-gonic/gin"
	"mini-alt/encoding"
	"net/http"
)

func respondS3Error(c *gin.Context, status int, code, message, bucket string) {
	err := encoding.S3Error{
		Code:       code,
		Message:    message,
		BucketName: bucket,
		RequestID:  "0000000000000000",
		HostID:     "local-s3-emulator",
	}
	c.XML(status, err)
}

type Error string

const (
	NoSuchBucket                  Error = "NoSuchBucket"
	BucketAlreadyExists                 = "BucketAlreadyExists"
	BucketIsNotEmpty                    = "BucketIsNotEmpty"
	FailedToCreateBucket                = "FailedToCreateBucket"
	FailedToDeleteBucket                = "FailedToDeleteBucket"
	FailedToDeleteBucketDirectory       = "FailedToDeleteBucketDirectory"
	FailedToDeleteObject                = "FailedToDeleteObject"
	FailedToDeleteObjectFile            = "FailedToDeleteObjectFile"
	NoSuchKey                           = "NoSuchKey"
	NoSuchMetadata                      = "NoSuchMetadata"
	EncryptionTypeMismatch              = "EncryptionTypeMismatch"
	InvalidRequest                      = "InvalidRequest"
	InvalidWriteOffset                  = "InvalidWriteOffset"
	TooManyParts                        = "TooManyParts"
	PreconditionFailed                  = "PreconditionFailed"
	InvalidSourceKey                    = "InvalidSourceKey"
	InvalidDestinationKey               = "InvalidDestinationKey"
	InternalServerError                 = "InternalServerError"
)

func HandleError(c *gin.Context, code Error, bucket string) {
	switch code {
	case NoSuchBucket:
		respondS3Error(c, http.StatusNotFound, "NoSuchBucket", "Bucket does not exist.", bucket)
	case BucketAlreadyExists:
		respondS3Error(c, http.StatusConflict, "BucketAlreadyExists", "The requested bucket name is not available.", bucket)
	case BucketIsNotEmpty:
		respondS3Error(c, http.StatusNotFound, "BucketIsNotEmpty", "Bucket is not empty.", bucket)
	case FailedToCreateBucket:
		respondS3Error(c, http.StatusInternalServerError, "FailedToCreateBucket", "Could not create storage bucket.", bucket)
	case FailedToDeleteBucket:
		respondS3Error(c, http.StatusInternalServerError, "FailedToDeleteBucket", "Could not delete bucket.", bucket)
	case FailedToDeleteBucketDirectory:
		respondS3Error(c, http.StatusInternalServerError, "FailedToDeleteBucketDirectory", "Could not delete bucket directory.", bucket)
	case FailedToDeleteObject:
		respondS3Error(c, http.StatusInternalServerError, "FailedToDeleteObject", "Could not delete object.", bucket)
	case FailedToDeleteObjectFile:
		respondS3Error(c, http.StatusInternalServerError, "FailedToDeleteObjectFile", "Could not delete object file.", bucket)
	case NoSuchKey:
		respondS3Error(c, http.StatusNotFound, "NoSuchKey", "Object not found.", bucket)
	case NoSuchMetadata:
		respondS3Error(c, http.StatusInternalServerError, "NoSuchMetadata", "Metadata not found.", bucket)
	case EncryptionTypeMismatch:
		respondS3Error(c, http.StatusBadRequest, "EncryptionTypeMismatch", "The existing object was created with a different encryption type. Subsequent write requests must include the appropriate encryption parameters in the request or while creating the session.", bucket)
	case InvalidRequest:
		respondS3Error(c, http.StatusBadRequest, "InvalidRequest", "Could not put object", bucket)
	case InvalidWriteOffset:
		respondS3Error(c, http.StatusBadRequest, "InvalidWriteOffset", "The write offset value that you specified does not match the current object size.", bucket)
	case TooManyParts:
		respondS3Error(c, http.StatusBadRequest, "TooManyParts", "You have attempted to add more parts than the maximum of 10000 that are allowed for this object. You can use the CopyObject operation to copy this object to another and then add more data to the newly copied object.", bucket)
	case PreconditionFailed:
		respondS3Error(c, http.StatusPreconditionFailed, "PreconditionFailed", "At least one of the preconditions you specified did not hold.", bucket)
	case InvalidSourceKey:
		respondS3Error(c, http.StatusBadRequest, "InvalidSourceKey", "The source key you specified is not valid.", bucket)
	case InvalidDestinationKey:
		respondS3Error(c, http.StatusBadRequest, "InvalidDestinationKey", "The destination key you specified is not valid.", bucket)
	case InternalServerError:
		respondS3Error(c, http.StatusInternalServerError, "InternalServerError", "Internal server error.", bucket)
	}
}
