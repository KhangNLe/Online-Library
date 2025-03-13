package loginsignup

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func UserLogIn(c *gin.Context, db *sql.DB) {
	var user Login
	user.UserName = c.PostForm("userid")
	user.Password = c.PostForm("password")

	hashPass, err := HashPass(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
	}

	c.JSON(200, gin.H{
		"message":  "data received",
		"username": user.UserName,
		"pass":     hashPass,
	})
}
