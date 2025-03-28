package main

import (
	"book/author"
	"book/htmxSwap"
	"book/login-signup"
	"book/search"
	user "book/user/add"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	projectPath = "Desktop/Projects/HTML/OnlineLibrary"
	clientPath  = "client-side"
	serverPath  = "server"
	dbName      = "library.db"
)

var key string

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
	gin.SetMode(gin.DebugMode)
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
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret-key"))
	if err != nil {
		log.Printf("Could not connect to redis. Error: %s", err)
	}
	r.Use(sessions.Sessions("mysesh", store))

	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		if session == nil {
			log.Printf("Could not conenct to session, Error: %s", err)
		}
		c.Next()
	})
	r.Static("/static", filepath.Join(frontFile, "static"))

	private := r.Group("/")
	private.Use(loginRequired())
	{
		private.GET("/my-books", func(ctx *gin.Context) {

		})
		private.GET("/user", func(c *gin.Context) {

		})
		private.GET("/logout", func(c *gin.Context) {
			session := sessions.Default(c)
			session.Delete("user_id")
			session.Set("authenticated", false)
			if err := session.Save(); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.Header("HX-Redirect", "/")
			c.Status(http.StatusOK)
		})
		private.GET("/wantToRead", func(c *gin.Context) {
			session := sessions.Default(c)
			userId, ok := session.Get("user_id").(string)
			if !ok {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			user.AddingToLibrary(userId, c, db, 0)
		})
		private.GET("/reading", func(c *gin.Context) {
			session := sessions.Default(c)
			userId, ok := session.Get("user_id").(string)
			if !ok {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			user.AddingToLibrary(userId, c, db, 1)
		})
		private.GET("/alreadyRead", func(c *gin.Context) {
			session := sessions.Default(c)
			userId, ok := session.Get("user_id").(string)
			if !ok {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			user.AddingToLibrary(userId, c, db, 69)
		})
	}
	r.LoadHTMLGlob(htmlFile)

	r.GET("/", func(c *gin.Context) { //TODO adding stuffs to the home page
		session := sessions.Default(c)
		auth := session.Get("authenticated")
		if auth == nil || auth.(bool) != true {
			c.HTML(http.StatusOK, "index.html", nil)
		} else {
			c.HTML(http.StatusOK, "user_logged_in.html", nil)
		}
	})

	r.GET("/log-in", htmxswap.LoginButton)

	r.GET("/search", func(c *gin.Context) {
		htmxRequest := c.GetHeader("HX-Request") == "true"
		if htmxRequest {
			search.SearchPage(c)
		} else {
			c.HTML(http.StatusOK, "search.html", nil)
		}
	})

	r.POST("/book-search", search.DisplaySearch)
	r.GET("/book-search", search.DisplaySearch)

	r.POST("/book", func(c *gin.Context) {
		search.BookDetail(c, db)
	})

	r.POST("/author", func(c *gin.Context) {
		author.GetAuthor(c, db)
	})

	r.GET("/recommend", func(c *gin.Context) {

	})

	r.GET("/sign-up", func(c *gin.Context) {
		htmxswap.SignUpBtn(c)
	})

	r.GET("/about", htmxswap.AboutPage)

	r.POST("/register", func(c *gin.Context) {
		loginsignup.RegisterUser(c, db)
	})

	r.POST("/user-login", func(c *gin.Context) {
		userNm, err := loginsignup.UserLogIn(c, db)
		if err != nil {
			log.Printf("Something wrong with the login process. Error: %s", err)
		} else {
			session := sessions.Default(c)
			session.Set("user_id", userNm)
			session.Set("authenticated", true)
			if err := session.Save(); err != nil {
				c.Header("Content-Type", "text/html")
				c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Session Error: Could not login into the account at the moment. Please try again later.
            </p> 
            `)

			} else {
				c.Header("HX-Redirect", "/")
				c.Status(http.StatusOK)
			}
		}

	})

	return r
}

func loginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		auth := session.Get("authenticated")

		if auth == nil || auth.(bool) != true {
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/user-login")
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.Redirect(http.StatusFound, "/")
				c.Abort()
			}
			return
		}
		c.Next()
	}
}
