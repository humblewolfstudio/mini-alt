package middlewares

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"mini-alt/handlers/api"
)

func BucketAuthentication(h *api.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, ok := c.Get("presignedAuth"); ok && val == true {
			c.Next()
			return
		}

		accessKey, exists := c.Get("accessKey")
		if !exists {
			respondS3Error(c, "InternalServerError", "User accessKey not found in context")
			return
		}

		user, err := h.Store.GetUserByAccessKey(accessKey.(string))
		if err != nil {
			respondS3Error(c, "InternalServerError", err.Error())
			return
		}
		c.Set("user", user)

		bucketName := c.Param("bucket")

		if len(bucketName) == 0 {
			c.Next()
			return
		}

		bucket, err := h.Store.GetBucket(bucketName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.Next()
				return
			}

			respondS3Error(c, "InternalServerError", err.Error())
			return
		}

		if user.Admin {
			c.Next()
			return
		}

		if bucket.Owner != user.Id {
			respondS3Error(c, "InternalServerError", "You do not own this bucket")
			return
		}

		c.Next()
	}
}
