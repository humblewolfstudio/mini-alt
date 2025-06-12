package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage"
	"mini-alt/utils"
	"net/http"
	"strconv"
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
	Expires            string `header:"Expires"` // RFC1123 or RFC3339 format
	IfMatch            string `header:"If-Match"`
	IfNoneMatch        string `header:"If-None-Match"`
}

type PutObjectErrors string

const (
	EncryptionTypeMismatch PutObjectErrors = "EncryptionTypeMismatch"
	InvalidRequest                         = "InvalidRequest"
	InvalidWriteOffset                     = "InvalidWriteOffset"
	TooManyParts                           = "TooManyParts"
)

// PutObject receives the bucket name, the object key and the object and persists it.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (h *Handler) PutObject(c *gin.Context, bucket, object string) {
	req, err := BindPutObjectRequest(c)

	if err != nil {
		HandleError(c, InvalidRequest, bucket, "Could not parse headers")
		return
	}

	path, err := storage.CreateObjectFilePath(bucket, object)
	if err != nil {
		HandleError(c, InvalidRequest, bucket, "Could not create object path")
		return
	}

	written, err := storage.CreateObject(path, c.Request.Body)
	if err != nil {
		HandleError(c, InvalidRequest, bucket, "Could not create object")
		return
	}

	h.Store.PutObject(bucket, object, written)
	c.Status(http.StatusOK)
}

func BindPutObjectRequest(c *gin.Context) (*PutObjectRequest, error) {
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
	}, nil
}

func HandleError(c *gin.Context, error PutObjectErrors, bucket, msg string) {
	switch error {
	case EncryptionTypeMismatch:
		utils.RespondS3Error(c, http.StatusBadRequest, "EncryptionTypeMismatch", "The existing object was created with a different encryption type. Subsequent write requests must include the appropriate encryption parameters in the request or while creating the session.", bucket)
	case InvalidRequest:
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidRequest", msg, bucket)
	case InvalidWriteOffset:
		utils.RespondS3Error(c, http.StatusBadRequest, "InvalidWriteOffset", "The write offset value that you specified does not match the current object size.", bucket)
	case TooManyParts:
		utils.RespondS3Error(c, http.StatusBadRequest, "TooManyParts", "ou have attempted to add more parts than the maximum of 10000 that are allowed for this object. You can use the CopyObject operation to copy this object to another and then add more data to the newly copied object.", bucket)
	}
}
