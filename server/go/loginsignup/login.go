package loginsignup

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Login struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func UserLogIn(c *gin.Context, db *sqlx.DB) (string, error) {
	var user Login
	user.UserName = c.PostForm("userid")
	user.Password = c.PostForm("password")

	resp, err := db.Query("SELECT * FROM User WHERE user_name=?", user.UserName)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Database Error: Could not login into the account at the moment. Please try again later.
            </p>
        `)
		return "", err
	}

	if !resp.Next() {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusBadRequest, `
            <br><p style="color: red; font-size: 14px;">
                Could not find any user with said username.
            </p>
        `)
		return "", err
	}
	resp.Close()

	resp, err = db.Query("SELECT pass_hash FROM User WHERE user_name=?", user.UserName)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Database Error: Could not login into your account at the moment. Please try again later.
            </p>
        `)
		return "", err
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
		return "", err
	}

	ans, err := db.Query(`SELECT user_id FROM User WHERE user_name = ?`, user.UserName)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <br><p style="color: red; font-size: 14px;">
                Database Error: Could not login into the account at the moment. Please try again later.
            </p>
        `)
		return "", err
	}
	defer ans.Close()

	if ans.Next() {
		var id int
		err = ans.Scan(&id)
		if err != nil {
			c.Header("Content-Type", "text/html")
			c.String(http.StatusInternalServerError, `
                <br><p style="color: red; font-size: 14px;">
                    Database Error: Could not login into the account at the moment. Please try again later.
                </p>
            `)
			return "", err
		}
		return strconv.Itoa(id), nil
	}
	return "", nil
}
