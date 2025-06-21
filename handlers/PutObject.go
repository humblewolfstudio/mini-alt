package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mini-alt/storage"
	"mini-alt/types"
	"mini-alt/utils"
	"net/http"
	"strconv"
)

type PutObjectErrors string

const (
	EncryptionTypeMismatch PutObjectErrors = "EncryptionTypeMismatch"
	InvalidRequest                         = "InvalidRequest"
	InvalidWriteOffset                     = "InvalidWriteOffset"
	TooManyParts                           = "TooManyParts"
	PreconditionFailed                     = "PreconditionFailed"
)

type PutObjectRequest struct {
	ACL                string `header:"x-amz-acl"`
	CacheControl       string `header:"Cache-Control"`
	ContentDisposition string `header:"Content-Disposition"`
	ContentEncoding    string `header:"Content-Encoding"`
	ContentLanguage    string `header:"Content-Language"`
	ContentLength      int64  `header:"Content-Length"`
	ContentMD5         string `header:"Content-MD5"`
	ContentType        string `header:"Content-Type"`
	Expires            string `header:"Expires"`
	IfMatch            string `header:"If-Match"`
	IfNoneMatch        string `header:"If-None-Match"`
}

func (req *PutObjectRequest) ToMetadata() types.Metadata {
	return types.Metadata{
		ACL:                req.ACL,
		CacheControl:       req.CacheControl,
		ContentDisposition: req.ContentDisposition,
		ContentEncoding:    req.ContentEncoding,
		ContentLanguage:    req.ContentLanguage,
		ContentLength:      req.ContentLength,
		ContentMD5:         req.ContentMD5,
		ContentType:        req.ContentType,
		Expires:            req.Expires,
	}
}

func BindPutObjectRequest(c *gin.Context) *PutObjectRequest {
	contentLength, _ := strconv.ParseInt(c.GetHeader("Content-Length"), 10, 64)
	return &PutObjectRequest{
		ACL:                c.GetHeader("x-amz-acl"),
		CacheControl:       c.GetHeader("Cache-Control"),
		ContentDisposition: c.GetHeader("Content-Disposition"),
		ContentEncoding:    c.GetHeader("Content-Encoding"),
		ContentLanguage:    c.GetHeader("Content-Language"),
		ContentLength:      contentLength,
		ContentMD5:         c.GetHeader("Content-MD5"),
		ContentType:        c.GetHeader("Content-Type"),
		Expires:            c.GetHeader("Expires"),
		IfMatch:            c.GetHeader("If-Match"),
		IfNoneMatch:        c.GetHeader("If-None-Match"),
	}
}

// PutObject receives the bucket name, the object key and the object and persists it.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (h *Handler) PutObject(c *gin.Context, bucketName, objectKey string) {
	_, err := h.Store.GetBucket(bucketName)
	if err != nil {
		_ = h.Store.PutBucket(bucketName)
	}

	path, err := storage.CreateObjectFilePath(bucketName, objectKey)
	if err != nil {
		HandleError(c, InvalidRequest, bucketName, "Could not create objectKey path")
		return
	}

	putObjectRequest := BindPutObjectRequest(c)

	if putObjectRequest.IfMatch != "" {
		existingObject, err := h.Store.GetObject(bucketName, objectKey)
		if err == nil {
			etag, err := storage.GetMD5Base64(existingObject.Key)
			if err == nil && etag != putObjectRequest.IfMatch {
				c.Header("ETag", etag)
				HandleError(c, PreconditionFailed, bucketName, "At least one of the preconditions you specified did not hold.")
				return
			}
		}
	}

	if putObjectRequest.IfNoneMatch != "" {
		existingObject, err := h.Store.GetObject(bucketName, objectKey)
		if err == nil {
			etag, err := storage.GetMD5Base64(existingObject.Key)
			if err == nil && etag == putObjectRequest.IfNoneMatch {
				c.Header("ETag", etag)
				HandleError(c, PreconditionFailed, bucketName, "An object already exists with the same ETag.")
				return
			}
		}
	}

	written, err := storage.CreateObject(path, c.Request.Body)
	if err != nil {
		HandleError(c, InvalidRequest, bucketName, "Could not write object")
		return
	}

	object, err := h.Store.PutObject(bucketName, objectKey, written)
	if err != nil {
		HandleError(c, InvalidRequest, bucketName, "Could not create object")
		return
	}

	md5, err := storage.GetMD5Base64(path)
	if err == nil {
		if md5 == putObjectRequest.ContentMD5 {
			fmt.Printf("Object integrity verified for %s", path)
		} else {
			fmt.Printf("Object integrity could not be verified for %s", path)
		}
	}

	err = h.Store.PutMetadata(object.Id, putObjectRequest.ToMetadata())
	if err != nil {
		HandleError(c, InvalidRequest, bucketName, "Could not create metadata")
		return
	}

	c.Header("ETag", md5)
	c.Status(http.StatusOK)
}

func HandleError(c *gin.Context, err PutObjectErrors, bucket, msg string) {
	switch err {
	case EncryptionTypeMismatch:
		utils.RespondS3Error(c, http.StatusBadRequest, "EncryptionTypeMismatch", "The existing object was created with a different encryption type. Subsequent write requests must include the appropriate encryption parameters in the request or while creating the session.", bucket)
	case InvalidRequest:
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidRequest", msg, bucket)
	case InvalidWriteOffset:
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidWriteOffset", "The write offset value that you specified does not match the current object size.", bucket)
	case TooManyParts:
		utils.RespondS3Error(c, http.StatusBadRequest, "TooManyParts", "You have attempted to add more parts than the maximum of 10000 that are allowed for this object. You can use the CopyObject operation to copy this object to another and then add more data to the newly copied object.", bucket)
	case PreconditionFailed:
		utils.RespondS3Error(c, http.StatusPreconditionFailed, "PreconditionFailed", msg, bucket)
	}
}
