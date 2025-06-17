package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/fs"
	"log"
	"mini-alt/router"
	"mini-alt/storage"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func main() {
	go startApiServer()
	startWebServer()
}

func startApiServer() {
	store, err := storage.NewSQLiteStore("./mini-alt.sqlite")
	if err != nil {
		log.Fatal(err)
		return
	}

	r := router.SetupRouter(store)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func startWebServer() {
	r := gin.Default()

	apiProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	})
	r.Any("/api/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		apiProxy.ServeHTTP(c.Writer, c.Request)
	})

	fsys, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(fsys))

	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.Request.URL.Path = "/assets" + c.Param("filepath")
		log.Printf("Serving asset: %s", c.Request.URL.Path)
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/vite.svg", func(c *gin.Context) {
		c.Request.URL.Path = "/vite.svg"
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/", serveIndex(fsys))
	r.NoRoute(serveIndex(fsys))

	if err := r.Run(":9001"); err != nil {
		log.Fatal(err)
	}
}

func serveIndex(fsys fs.FS) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := fsys.Open("index.html")
		if err != nil {
			log.Printf("Error opening index.html: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		defer func(file fs.File) {
			err := file.Close()
			if err != nil {
				println("Error closing file: %v", err)
			}
		}(file)

		content, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading index.html: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	}
}
