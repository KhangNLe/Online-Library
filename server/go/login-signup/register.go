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

	query, err := db.Beginx()
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusNotFound, `
            <p style="color:red; font-size: 14px;">Error: Database Service error, please try again later.</p>
            `)
		log.Printf("Data service error: %s", err)
		return
	}
	defer func() {
		if err != nil {
			query.Rollback()
		}
	}()

	existenUser, err := query.Query("SELECT user_id FROM USER WHERE user_name=(?)", userSignup.UserName)
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

	if _, err := query.Exec("INSERT INTO USER(user_name, pass_hash) VALUES (?, ?)", userSignup.UserName, hashPass); err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusNotAcceptable, `
            <p style="color:red; font-size: 14px;">Database Service error: We could not register you at the current moment, please try again later.</p>
            `)
		log.Printf("Data service error: %s", err)
		return
	}

	resp, err := query.Query(`SELECT user_id FROM User WHERE user_name = ?`, userSignup.UserName)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Database Error: Could not register into the account at the moment. Please try again later.
            </p>
        `)
		log.Printf("Error getting user_id. Error: %s", err)
		return
	}
	defer resp.Close()

	if resp.Next() {
		var id int
		err = resp.Scan(&id)
		if err != nil {
			c.Header("Content-Type", "text/html")
			c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Database Error: Could not register into the account at the moment. Please try again later.
            </p>
        `)
			log.Printf("Error scanning id. Error: %s", err)
			return
		}
		_, err = query.Exec("INSERT INTO User_library(user_id) VALUES (?)", id)
		if err != nil {
			c.Header("Content-Type", "text/html")
			c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Database Error: Could not register into the account at the moment. Please try again later.
            </p>
        `)
			log.Printf("Error insert user_id into User_Library. Error: %s", err)
			return
		}

		if err = query.Commit(); err != nil {
			c.Header("Content-Type", "text/html")
			c.String(http.StatusNotFound, `
            <p style="color:red; font-size: 14px;">Error: Database Service error, please try again later.</p>
            `)
			log.Printf("Data service error: %s", err)
			return
		}
		c.Header("HX-Redirect", "/")
		c.Status(http.StatusOK)
	} else {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Database Error: Could not register into the account at the moment. Please try again later.
            </p>
        `)
		return
	}
}
