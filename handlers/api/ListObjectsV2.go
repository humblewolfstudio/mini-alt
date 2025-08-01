package api

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
	objects, err := h.Store.ListObjects(bucket)
	if err != nil {
		utils.RespondS3Error(c, http.StatusNotFound, "NoSuchBucket", "The specified bucket does not exist.", bucket)
		return
	}

	params := c.Request.URL.Query()
	var maxKeys = MaxObjects
	var startAfter = ""
	var prefix = ""
	var delimiter = ""

	// Parse max-keys parameter
	if params.Has("max-keys") {
		if parsedMaxKeys, err := strconv.Atoi(params.Get("max-keys")); err == nil {
			maxKeys = parsedMaxKeys
		}
	}

	// Parse start-after parameter
	if params.Has("start-after") {
		startAfter = strings.TrimSpace(params.Get("start-after"))
	}

	// Parse prefix parameter
	if params.Has("prefix") {
		prefix = strings.TrimSpace(params.Get("prefix"))
		prefix = strings.TrimPrefix(prefix, "/")
	}

	// Parse delimiter parameter (typically "/")
	if params.Has("delimiter") {
		delimiter = params.Get("delimiter")
	}

	var xmlListBucketResult encoding.ListBucketResult
	xmlListBucketResult.Name = bucket
	xmlListBucketResult.MaxKeys = maxKeys
	xmlListBucketResult.StartAfter = startAfter

	var xmlContents []encoding.Content
	var commonPrefixes []encoding.CommonPrefix

	keyCount := 0
	isTruncated := false
	foundStartAfter := startAfter == ""
	seenPrefixes := make(map[string]bool)

	for _, object := range objects {
		if !foundStartAfter {
			if object.Key == startAfter {
				foundStartAfter = true
			}
			continue
		}

		if keyCount >= maxKeys {
			isTruncated = true
			break
		}

		objectKey := strings.TrimPrefix(object.Key, "/")

		if prefix != "" && !strings.HasPrefix(objectKey, prefix) {
			continue
		}

		if delimiter != "" {
			remainingPart := strings.TrimPrefix(objectKey, prefix)
			delimiterPos := strings.Index(remainingPart, delimiter)

			if delimiterPos >= 0 {
				commonPrefix := prefix + remainingPart[:delimiterPos+len(delimiter)]

				if !seenPrefixes[commonPrefix] {
					commonPrefixes = append(commonPrefixes, encoding.CommonPrefix{
						Prefix: commonPrefix,
					})
					seenPrefixes[commonPrefix] = true
					keyCount++
				}
				continue
			}
		}

		xmlContents = append(xmlContents, encoding.Content{
			Key:          object.Key,
			LastModified: object.LastModified,
			Size:         object.Size,
		})
		keyCount++
	}

	xmlListBucketResult.Contents = xmlContents
	xmlListBucketResult.CommonPrefixes = commonPrefixes
	xmlListBucketResult.IsTruncated = isTruncated
	xmlListBucketResult.KeyCount = keyCount

	if delimiter != "" {
		xmlListBucketResult.Delimiter = delimiter
	}

	c.XML(http.StatusOK, xmlListBucketResult)
}
