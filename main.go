package main

import (
	"embed"
	"flag"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/fs"
	"log"
	"mini-alt/router"
	"mini-alt/storage"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func main() {
	noWeb := flag.Bool("no-web", false, "Disable web interface")
	flag.Parse()

	if *noWeb {
		startApiServer()
		println("Starting without web interface")
	} else {
		go startApiServer()
		startWebServer()
	}
}

func startApiServer() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	appDir := filepath.Dir(exe)

	dbPath := filepath.Join(appDir, "mini-alt.sqlite")
	store, err := storage.NewSQLiteStore(dbPath)

	if err != nil {
		log.Fatal(err)
		return
	}

	r := router.SetupAPIRouter(store)

	if err := r.Run(":9000"); err != nil {
		log.Fatal(err)
	}
}

func startWebServer() {
	r := router.SetupWebRouter()

	fsys, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(fsys))

	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.Request.URL.Path = "/assets" + c.Param("filepath")
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/vite.svg", func(c *gin.Context) {
		c.Request.URL.Path = "/vite.svg"
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/", serveIndex(fsys))

	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api/") {
			serveIndex(fsys)(c)
			return
		}
		c.JSON(404, gin.H{"error": "Not found"})
	})

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
