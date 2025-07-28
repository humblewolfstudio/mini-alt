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
func (h *ApiHandler) ListObjectsV2(c *gin.Context, bucket string) {
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
		parsedMaxKeys, err := strconv.Atoi(params.Get("max-keys"))
		if err == nil {
			maxKeys = parsedMaxKeys
		}
	}

	// Parse start-after parameter
	if params.Has("start-after") {
		startAfter = params.Get("start-after")
	}

	// Parse prefix parameter
	if params.Has("prefix") {
		prefix = params.Get("prefix")
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

		if len(prefix) > 0 && !strings.HasPrefix(object.Key, prefix) {
			continue
		}

		if delimiter != "" {
			remainingPart := strings.TrimPrefix(object.Key, prefix)

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
				continue // Skip adding to Contents as it's in a subdirectory
			}
		}

		// Add to Contents only if:
		// 1. No prefix is specified, or
		// 2. The object is exactly in the current prefix level (no additional delimiters)
		if len(prefix) == 0 || !strings.Contains(strings.TrimPrefix(object.Key, prefix), delimiter) {
			xmlContents = append(xmlContents, encoding.Content{
				Key:          object.Key,
				LastModified: object.LastModified,
				Size:         object.Size,
			})
			keyCount++
		}
	}

	// Set the response fields
	xmlListBucketResult.Contents = xmlContents
	xmlListBucketResult.CommonPrefixes = commonPrefixes
	xmlListBucketResult.IsTruncated = isTruncated
	xmlListBucketResult.KeyCount = keyCount

	if delimiter != "" {
		xmlListBucketResult.Delimiter = delimiter
	}

	c.XML(http.StatusOK, xmlListBucketResult)
}
