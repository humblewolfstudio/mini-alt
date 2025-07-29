package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mini-alt/handlers/api"
	"mini-alt/handlers/web"
	"mini-alt/middlewares"
	"mini-alt/storage"
)

func SetupAPIRouter(store storage.Store) *gin.Engine {
	r := gin.New()
	r.Use(customLogger("API-SERVER"))

	h := api.ApiHandler{Store: store}

	r.Use(middlewares.APIAuthenticationMiddleware(&h))

	r.GET("/", h.ListBuckets)
	r.PUT("/:bucket/*object", h.PutObjectOrBucket)
	r.GET("/:bucket/*object", h.GetObjectOrList)
	r.HEAD("/:bucket/*object", h.HeadObjectOrBucket)
	r.DELETE("/:bucket/*object", h.DeleteObjectOrBucket)

	return r
}

func SetupWebRouter(store storage.Store) *gin.Engine {
	r := gin.New()
	r.Use(customLogger("WEB-SERVER"))

	h := web.WebHandler{Store: store}

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/buckets", h.ListBuckets)
		apiGroup.POST("/buckets", web.PutBucket)
		apiGroup.GET("/files/list", web.ListFiles)
		apiGroup.GET("/files/list-folders", web.ListFolders)
		apiGroup.GET("/files/download", web.DownloadFile)
		apiGroup.POST("/files/upload", web.UploadFiles)
		apiGroup.POST("/files/create-folder", web.CreateFolder)
		apiGroup.POST("/files/delete", web.DeleteFile)
		apiGroup.PUT("/files/rename", web.RenameFile)
		apiGroup.PUT("/files/move", web.MoveFile)

		apiGroup.GET("/credentials", h.ListCredentials)
		apiGroup.POST("/credentials", h.CreateCredentials)
		apiGroup.POST("/credentials/delete", h.DeleteCredentials)
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
