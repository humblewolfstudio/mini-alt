package main

import (
	"github.com/gin-gonic/gin"
	"mini-alt/storage"
	"net/http"
	"os"
)

var store = storage.NewInMemoryStore()

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.PUT("/buckets/:bucket", func(c *gin.Context) {
		bucket := c.Param("bucket")

		if err := store.CreateBucket(bucket); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := os.MkdirAll("uploads/"+bucket, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
