package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mini-alt/handlers"
	"mini-alt/handlers/web"
	"mini-alt/storage"
)

func SetupAPIRouter(store storage.Store) *gin.Engine {
	r := gin.New()
	r.Use(customLogger("API-SERVER"))

	h := handlers.Handler{Store: store}

	r.GET("/", h.ListBuckets)
	r.PUT("/:bucket/*object", h.PutObjectOrBucket)
	r.GET("/:bucket/*object", h.GetObjectOrList)
	r.HEAD("/:bucket/*object", h.HeadObject)
	r.DELETE("/:bucket/*object", h.DeleteObjectOrBucket)

	return r
}

func SetupWebRouter() *gin.Engine {
	r := gin.New()
	r.Use(customLogger("WEB-SERVER"))

	api := r.Group("/api")
	{
		api.GET("/buckets", web.ListBuckets)
	}

	return r
}

func customLogger(serverName string) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s - [%s] \"%s %s\" %d %s %s\n",
			serverName,
			param.ClientIP,
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	})
}
