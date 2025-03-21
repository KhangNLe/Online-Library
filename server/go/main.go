package main

import (
	"book/homepage"
	"book/htmxSwap"
	"book/login-signup"
	"book/search"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	projectPath = "Desktop/Projects/HTML/OnlineLibrary"
	clientPath  = "client-side"
	serverPath  = "server"
	dbName      = "library.db"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not open local home directory. Error: %s", err)
	}

	frontFile := filepath.Join(home, projectPath, clientPath)
	htmlFile := filepath.Join(frontFile, "*.html")
	dbFile := filepath.Join(home, projectPath, serverPath, dbName)

	db, err := connectDB(dbFile)
	if err != nil {
		log.Fatalf("Could not connect to database. Error: %s", err)
	}
	defer db.Close()

	r := setupRouter(frontFile, htmlFile, db)

	r.Run("localhost:6969")
}

func connectDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupRouter(frontFile, htmlFile string, db *sqlx.DB) *gin.Engine {
	r := gin.Default()
	r.Static("/static", filepath.Join(frontFile, "static"))

	r.LoadHTMLGlob(htmlFile)

	r.GET("/", func(c *gin.Context) { //TODO adding stuffs to the home page
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/log-in", htmxswap.LoginButton)

	r.GET("/search", func(c *gin.Context) {
		htmxRequest := c.GetHeader("HX-Request") == "true"

		if htmxRequest {
			search.SearchPage(c)
		} else {
			c.HTML(http.StatusOK, "search.html", gin.H{"message": "Search"})
		}
	})

	r.POST("/book-search", search.DisplaySearch)
	r.GET("/book-search", search.DisplaySearch)

	r.GET("/my-books", func(ctx *gin.Context) {

	})

	r.POST("/book", func(c *gin.Context) {
		search.BookDetail(c, db)
	})

	r.GET("/recommend", homepage.Homepage)

	r.GET("/sign-up", func(c *gin.Context) {
		htmxswap.SignUpBtn(c)
	})

	r.GET("/about", htmxswap.AboutPage)

	r.POST("/register", func(c *gin.Context) {
		loginsignup.RegisterUser(c, db)
	})

	r.POST("/user-login", func(c *gin.Context) {
		loginsignup.UserLogIn(c, db)
	})

	r.GET("/user", func(c *gin.Context) {

	})

	return r
}
