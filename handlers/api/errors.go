package api

import (
	"github.com/gin-gonic/gin"
	"mini-alt/utils"
	"net/http"
)

type Error string

const (
	NoSuchBucket                  Error = "NoSuchBucket"
	BucketAlreadyExists                 = "BucketAlreadyExists"
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
)

func handleError(c *gin.Context, code Error, bucket string) {
	switch code {
	case NoSuchBucket:
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchBucket", "Bucket does not exist.", bucket)
	case BucketAlreadyExists:
		utils.RespondS3Error(c, http.StatusConflict, "BucketAlreadyExists", "The requested bucket name is not available.", bucket)
	case FailedToCreateBucket:
		utils.RespondS3Error(c, http.StatusInternalServerError, "FailedToCreateBucket", "Could not create storage bucket.", bucket)
	case FailedToDeleteBucket:
		utils.RespondS3Error(c, http.StatusInternalServerError, "FailedToDeleteBucket", "Could not delete bucket.", bucket)
	case FailedToDeleteBucketDirectory:
		utils.RespondS3Error(c, http.StatusInternalServerError, "FailedToDeleteBucketDirectory", "Could not delete bucket directory.", bucket)
	case FailedToDeleteObject:
		utils.RespondS3Error(c, http.StatusInternalServerError, "FailedToDeleteObject", "Could not delete object.", bucket)
	case FailedToDeleteObjectFile:
		utils.RespondS3Error(c, http.StatusInternalServerError, "FailedToDeleteObjectFile", "Could not delete object file.", bucket)
	case NoSuchKey:
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchKey", "Object not found.", bucket)
	case NoSuchMetadata:
		utils.RespondS3Error(c, http.StatusInternalServerError, "NoSuchMetadata", "Metadata not found.", bucket)
	}
}
