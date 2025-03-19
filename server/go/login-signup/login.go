package loginsignup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Login struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func UserLogIn(c *gin.Context, db *sqlx.DB) {
	var user Login
	user.UserName = c.PostForm("userid")
	user.Password = c.PostForm("password")

	resp, err := db.Query("SELECT * FROM USER WHERE user_name=?", user.UserName)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Database Error: Could not login into the account at the moment. Please try again later.
            </p>
        `)
		return
	}

	if !resp.Next() {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusBadRequest, `
            <br><p style="color: red; font-size: 14px;">
                Could not find any user with said username.
            </p>
        `)
		return
	}
	resp.Close()

	resp, err = db.Query("SELECT hash_pass FROM USER WHERE user_name=?", user.UserName)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Database Error: Could not login into your account at the moment. Please try again later.
            </p>
        `)
		return
	}
	defer resp.Close()

	var hashPass string
	for resp.Next() {
		resp.Scan(&hashPass)
	}

	match := ComparePass(hashPass, user.Password)

	if !match {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusBadRequest, `
            <br><p style="color: red; font-size: 14px;">
                Username or password did not match, please try again.
            </p>
        `)
		return
	}

	c.Header("HX-Redirect", "/user")
	c.Status(http.StatusOK)
}
