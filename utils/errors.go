package utils

import (
	"github.com/gin-gonic/gin"
	"mini-alt/encoding"
)

func RespondS3Error(c *gin.Context, status int, code, message, bucket string) {
	err := encoding.S3Error{
		Code:       code,
		Message:    message,
		BucketName: bucket,
		RequestID:  "0000000000000000",
		HostID:     "local-s3-emulator",
	}
	c.XML(status, err)
}
