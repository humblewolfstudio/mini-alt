package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"mini-alt/encoding"
	"mini-alt/storage"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var store = storage.NewInMemoryStore()

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		buckets := store.ListBuckets()
		var xmlBuckets []encoding.BucketXML

		for _, bucket := range buckets {
			xmlBuckets = append(xmlBuckets, encoding.BucketXML{
				Name:      bucket.Name,
				CreatedAt: bucket.CreatedAt.Format(time.RFC3339),
			})
		}

		res := encoding.ListAllMyBucketsResult{}
		res.Buckets.Bucket = xmlBuckets

		c.XML(http.StatusOK, res)
	})

	r.PUT("/:bucket/*object", func(c *gin.Context) {
		bucket := c.Param("bucket")
		object := c.Param("object")

		// If object is "/" or empty, treat as bucket creation
		if object == "/" || object == "" {
			handleCreateBucket(c, bucket)
			return
		}

		// Strip leading slash for object key
		object = object[1:]
		handleUploadObject(c, bucket, object)
	})

	r.GET("/:bucket/*object", func(c *gin.Context) {
		bucket := c.Param("bucket")
		object := c.Param("object")

		path := filepath.Join("uploads", bucket, object)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			respondS3Error(c, http.StatusConflict, "ObjectNotFound",
				"Object not found.",
				bucket)
			return
		}

		c.File(path)
	})

	if err := r.Run(":8080"); err != nil {
		return
	}
}

func handleCreateBucket(c *gin.Context, bucket string) {
	if err := store.CreateBucket(bucket); err != nil {
		respondS3Error(c, http.StatusConflict, "BucketAlreadyExists",
			"The requested bucket name is not available. The bucket namespace is shared by all users of the system. Select a different name and try again.",
			bucket)
		return
	}

	if err := os.MkdirAll("uploads/"+bucket, os.ModePerm); err != nil {
		respondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create storage directory.", bucket)
		return
	}

	c.Status(http.StatusCreated)
}

func handleUploadObject(c *gin.Context, bucket, object string) {
	dst := filepath.Join("uploads", bucket, object)
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		respondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create directory path.", bucket)
		return
	}

	file, err := os.Create(dst)
	if err != nil {
		println(err.Error())
		respondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not create file.", bucket)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	written, err := io.Copy(file, c.Request.Body)
	if err != nil {
		respondS3Error(c, http.StatusInternalServerError, "InternalError",
			"Could not write file.", bucket)
		return
	}

	store.PutObject(bucket, object, written)
	c.Status(http.StatusOK)
}

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
