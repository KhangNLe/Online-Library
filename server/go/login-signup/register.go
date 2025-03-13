package loginsignup

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func RegisterUser(c *gin.Context) {
	var userSignup Login
	userSignup.UserName = c.PostForm("userid")
	userSignup.Password = c.PostForm("password")

	hashPass, err := HashPass(userSignup.Password)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <p style="color:red; font-size: 14px;">Error: Unable to hash password</p>
            `)
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <p style="color:red; font-size: 14px;">Error: Backend Service error, please try again later.</p>
            `)
		return
	}

	dbDir := filepath.Join(homeDir, "Desktop/Projects/HTML/OnlineLibrary/server/library.db")
	db, err := sql.Open("sqlite3", dbDir)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <p style="color:red; font-size: 14px;">Error: Backend Service error, please try again later.</p>
            `)
		return
	}

	existenUser, err := db.Query("SELECT user_id FROM USER WHERE user_name=(?)", userSignup.UserName)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusNotFound, `
            <p style="color:red; font-size: 14px;">Error: Database Service error, please try again later.</p>
            `)
		return
	}
	defer existenUser.Close()

	if existenUser.Next() {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusBadRequest, `
            <p style="color:red; font-size: 14px;">Username is already taken. Please try another name.</p>
            `)
		return
	}

	defer db.Close()
	if _, err := db.Exec("INSERT INTO USER(user_name, hash_pass) VALUES (?, ?)", userSignup.UserName, hashPass); err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusNotAcceptable, `
            <p style="color:red; font-size: 14px;">Error: Database Service error: We could not register you at the current moment, please try again later.</p>
            `)
		return
	}

	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}
