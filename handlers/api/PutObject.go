package api

import (
	"github.com/gin-gonic/gin"
	"mime"
	"mime/multipart"
	"mini-alt/events"
	eventsTypes "mini-alt/events/types"
	"mini-alt/storage/disk"
	"mini-alt/types"
	"mini-alt/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
func (h *Handler) PutObject(c *gin.Context, bucket, objectKey string) {
	clientIp := utils.GetClientIP(c.Request)
	accessKey := c.GetString("accessKey")
	if accessKey == "" {
		utils.HandleError(c, utils.InternalServerError, "Access key not found")
		return
	}

	putObjectRequest := BindPutObjectRequest(c)

	if c.Request.MultipartForm != nil {
		fileHeader, err := c.FormFile("file")
		if err == nil {
			putObjectRequest = fillMissingMetadataFromFile(fileHeader, putObjectRequest)
		}
	} else {
		if putObjectRequest.ContentType == "" {
			ext := filepath.Ext(objectKey)
			putObjectRequest.ContentType = mime.TypeByExtension(ext)
		}
	}

	if putObjectRequest.IfMatch != "" {
		existingObject, err := h.Store.GetObject(bucket, objectKey)
		if err == nil {
			etag, err := disk.GetMD5Base64(bucket, existingObject.Key)
			if err == nil && etag != putObjectRequest.IfMatch {
				c.Header("ETag", etag)
				utils.HandleError(c, utils.PreconditionFailed, bucket)
				return
			}
		}
	}

	if putObjectRequest.IfNoneMatch != "" {
		existingObject, err := h.Store.GetObject(bucket, objectKey)
		if err == nil {
			etag, err := disk.GetMD5Base64(bucket, existingObject.Key)
			if err == nil && etag == putObjectRequest.IfNoneMatch {
				c.Header("ETag", etag)
				utils.HandleError(c, utils.PreconditionFailed, bucket)
				return
			}
		}
	}

	etag, e := h.Storage.PutObject(bucket, objectKey, c.Request.Body, putObjectRequest.ToMetadata())
	if e != "" {
		utils.HandleError(c, e, bucket)
		return
	}

	go events.HandleEventObject(h.Store, eventsTypes.EventPut, utils.ClearInput(bucket), utils.ClearInput(objectKey), etag, putObjectRequest.ContentLength, accessKey, clientIp)

	c.Header("ETag", etag)
	c.Status(http.StatusOK)
}

func fillMissingMetadataFromFile(fileHeader *multipart.FileHeader, req *PutObjectRequest) *PutObjectRequest {
	if req == nil {
		req = &PutObjectRequest{}
	}

	if req.ContentType == "" {
		ext := filepath.Ext(fileHeader.Filename)
		req.ContentType = mime.TypeByExtension(ext)
	}

	if req.ContentLength == 0 {
		req.ContentLength = fileHeader.Size
	}

	if req.ContentDisposition == "" {
		req.ContentDisposition = "inline; filename=\"" + fileHeader.Filename + "\""
	}

	if req.ContentEncoding == "" && strings.HasSuffix(fileHeader.Filename, ".gz") {
		req.ContentEncoding = "gzip"
	}

	return req
}
