package main

import (
	"book/htmxSwap"
	"book/login-signup"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not open local home directory. Error: %s", err)
	}

	frontFile := filepath.Join(home, "Desktop/Projects/HTML/OnlineLibrary/client-side/")
	htmlFile := filepath.Join(frontFile, "index.html")

	r.Static("/static", filepath.Join(frontFile, "static"))

	r.LoadHTMLFiles(htmlFile)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"message": "Hello there"})
	})

	r.GET("/log-in", func(c *gin.Context) {
		htmxswap.LoginButton(c)
	})

	r.GET("/sign-up", func(c *gin.Context) {
		htmxswap.SignUpBtn(c)
	})

	r.GET("/about", func(ctx *gin.Context) {
		htmxswap.AboutPage(ctx)
	})

	r.POST("/user-log-in", func(ctx *gin.Context) {
		loginsignup.UserLogIn(ctx)
	})
	r.Run("localhost:6969")
}
