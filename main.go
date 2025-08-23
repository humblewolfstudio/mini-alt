package main

import (
	"embed"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/fs"
	"log"
	"mini-alt/crons"
	"mini-alt/events"
	"mini-alt/jobs"
	"mini-alt/router"
	"mini-alt/storage/db"
	"mini-alt/storage/disk"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

// TODO add test flag that initializes the database with a test user (admin) and a test credentials that are always the same. This way, we can test it in the cloud
func main() {
	loadInitialData := flag.Bool("load-initial-data", false, "Load initial data")
	testData := flag.Bool("test", false, "Initializes the app with a test user/access key")
	flag.Parse()

	loadEnv()
	store := startDatabase()
	crons.StartupCronJobs(store)

	if *testData {
		jobs.LoadTestCredentials(store)
	}

	if *loadInitialData {
		jobs.LoadInitialData(store)
	}

	events.InitPool(10)

	go startApiServer(store)
	startWebServer(store)
}

func startDatabase() *db.Store {
	configDir, err := disk.GetAppConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	parentDir := filepath.Dir(configDir)
	println("Starting server in", parentDir)

	path := filepath.Join(configDir, "mini-alt.sqlite")
	store, err := db.NewStore(path)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return store
}

func loadEnv() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: Error loading .env file - relying on system environment variables")
		}
	}
}

func startApiServer(store *db.Store) {

	r := router.SetupAPIRouter(store)

	if err := r.Run(":9000"); err != nil {
		log.Fatal(err)
	}
}

func startWebServer(store *db.Store) {
	r := router.SetupWebRouter(store)

	fsys, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(fsys))

	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.Request.URL.Path = "/assets" + c.Param("filepath")
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/icons/*filepath", func(c *gin.Context) {
		c.Request.URL.Path = "/icons" + c.Param("filepath")
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
