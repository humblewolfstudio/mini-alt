package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mini-alt/handlers/api"
	"mini-alt/handlers/web"
	"mini-alt/middlewares"
	"mini-alt/storage/db"
)

func SetupAPIRouter(store *db.Store) *gin.Engine {
	r := gin.New()
	r.RedirectTrailingSlash = false
	r.Use(customLogger("API-SERVER"))

	h := api.Handler{Store: store}

	r.Use(middlewares.APIAuthenticationMiddleware(&h))

	r.GET("/", h.ListBuckets)
	r.PUT("/:bucket/*object", h.PutObjectOrBucket)
	r.GET("/:bucket", h.GetObjectOrList)
	r.GET("/:bucket/*object", h.GetObjectOrList)
	r.HEAD("/:bucket/*object", h.HeadObjectOrBucket)
	r.DELETE("/:bucket/*object", h.DeleteObjectOrBucket)

	return r
}

func SetupWebRouter(store *db.Store) *gin.Engine {
	r := gin.New()
	r.Use(customLogger("WEB-SERVER"))

	h := web.Handler{Store: store}
	r.POST("/api/users/login", h.LoginUser)

	apiGroup := r.Group("/api")
	apiGroup.Use(middlewares.WebAuthenticationMiddleware(&h))
	{
		apiGroup.GET("/buckets", h.ListBuckets)
		apiGroup.POST("/buckets", h.PutBucket)
		apiGroup.GET("/files/list", h.ListFiles)
		apiGroup.GET("/files/list-folders", h.ListFolders)
		apiGroup.GET("/files/download", h.DownloadFile)
		apiGroup.POST("/files/upload", h.UploadFiles)
		apiGroup.POST("/files/create-folder", h.CreateFolder)
		apiGroup.POST("/files/delete", h.DeleteFile)
		apiGroup.PUT("/files/rename", h.RenameFile)
		apiGroup.PUT("/files/move", h.MoveFile)

		apiGroup.GET("/credentials", h.ListCredentials)
		apiGroup.POST("/credentials", h.CreateCredentials)
		apiGroup.POST("/credentials/delete", h.DeleteCredentials)
		apiGroup.POST("/credentials/edit", h.CredentialsEdit)

		apiGroup.GET("/users/list", h.ListUsers)
		apiGroup.POST("/users/register", h.RegisterUser)
		apiGroup.POST("/users/delete", h.DeleteUser)
		apiGroup.GET("/users/authenticate", h.AuthenticateUser)

		apiGroup.GET("/users/logout", h.LogoutUser)
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
