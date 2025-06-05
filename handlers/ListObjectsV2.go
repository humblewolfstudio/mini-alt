package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-alt/encoding"
	"mini-alt/utils"
	"net/http"
	"strconv"
	"strings"
)

const MaxObjects = 1000

// ListObjectsV2 returns a list of x objects in a bucket.
// AWS Documentation: https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html
func (h *Handler) ListObjectsV2(c *gin.Context, bucket string) {
	objects := h.Store.ListObjects(bucket)
	if objects == nil {
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchBucket", "The specified bucket does not exist.", bucket)
		return
	}

	// Request parameters
	params := c.Request.URL.Query()
	var maxKeys = MaxObjects
	var startAfter = ""
	var prefix = ""

	// If max-keys exist, it has a value, and it's a correct integer we change the default maxKeys to the parameter one.
	if params.Has("max-keys") {
		parsedMaxKeys, err := strconv.Atoi(params.Get("max-keys"))
		if err == nil {
			maxKeys = parsedMaxKeys
		}
	}

	// start-after is where you want to start listing from.
	if params.Has("start-after") {
		startAfter = params.Get("start-after")
	}

	// prefix limits the response to keys that begin with the specified prefix.
	if params.Has("prefix") {
		prefix = params.Get("prefix")
	}

	// XML Parsing
	var xmlListBucketResult encoding.ListBucketResult

	xmlListBucketResult.Name = bucket
	xmlListBucketResult.MaxKeys = maxKeys
	xmlListBucketResult.StartAfter = startAfter

	// If startAfter is different from 0 (default) we return in the response the value
	if len(startAfter) > 0 {
		xmlListBucketResult.StartAfter = startAfter
	}

	var xmlContents []encoding.Content

	keyCount := 0
	isTruncated := false

	hasPrefix := len(prefix) > 0
	foundStartAfter := startAfter == ""

	for _, object := range objects {
		if !foundStartAfter {
			if object.ObjectKey == startAfter {
				foundStartAfter = true
			}
			continue
		}

		if keyCount >= maxKeys {
			isTruncated = true
			break
		}

		if hasPrefix && !strings.HasPrefix(object.ObjectKey, prefix) {
			continue
		}

		xmlContents = append(xmlContents, encoding.Content{
			Key:          object.ObjectKey,
			LastModified: object.LastModified,
			Size:         object.Size,
		})

		keyCount++
	}

	xmlListBucketResult.Contents = xmlContents
	xmlListBucketResult.IsTruncated = isTruncated
	xmlListBucketResult.KeyCount = keyCount

	c.XML(http.StatusOK, xmlListBucketResult)
}
