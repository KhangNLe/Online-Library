package loginsignup

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterUser(c *gin.Context, db *sqlx.DB) {
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

	existenUser, err := db.Query("SELECT user_id FROM USER WHERE user_name=(?)", userSignup.UserName)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusNotFound, `
            <p style="color:red; font-size: 14px;">Error: Database Service error, please try again later.</p>
            `)
		log.Printf("Data service error: %s", err)
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

	if _, err := db.Exec("INSERT INTO USER(user_name, pass_hash) VALUES (?, ?)", userSignup.UserName, hashPass); err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusNotAcceptable, `
            <p style="color:red; font-size: 14px;">Database Service error: We could not register you at the current moment, please try again later.</p>
            `)
		log.Printf("Data service error: %s", err)
		return
	}

	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}
