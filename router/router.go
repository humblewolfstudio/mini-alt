package router

import (
	"github.com/gin-gonic/gin"
	"mini-alt/handlers"
	"mini-alt/storage"
)

func SetupRouter(store storage.Store) *gin.Engine {
	r := gin.Default()
	h := handlers.Handler{Store: store}

	r.GET("/", h.ListBuckets)
	r.PUT("/:bucket/*object", h.PutObjectOrBucket)
	r.GET("/:bucket/*object", h.GetObjectOrList)

	return r
}
