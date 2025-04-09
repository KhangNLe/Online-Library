package user

import (
	"book/mybook"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func UpdateProfile(c *gin.Context, db *sqlx.DB, user string) {
	userId, err := strconv.Atoi(user)
	if err != nil {
		mybook.ErrorRespone(c, `

			`, http.StatusInternalServerError)
		log.Printf("User: %s cannot change it into int. Error: %s", err, err)
		return
	}

	fname := c.PostForm("fname")
	lname := c.PostForm("lname")
	email := c.PostForm("email")

	tx, err := db.Beginx()
	if err != nil {
		mybook.ErrorRespone(c, `

			`, http.StatusInternalServerError)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = updateProfile(tx, fname, lname, email, userId)
	if err != nil {
		mybook.ErrorRespone(c, `

			`, http.StatusInternalServerError)
		log.Printf("Could not update profile. Error: %s", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		mybook.ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Could not commit into db. Error: %s", err)
		return
	}
}

func updateProfile(tx *sqlx.Tx, fname, lname, email string,
	userId int) error {
	emailRegex := `(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|
		"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")
		@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|
		[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])
		|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|
		\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`

	regex := regexp.MustCompile(emailRegex)

	if !regex.MatchString(email) {
		return errors.New(fmt.Sprintf("Invalid email of %s", email))
	}

	_, err := tx.Exec(`INSERT INTO User_Profile 
			VALUES (?, ? ,?, ?)`, userId, fname, lname, email)
	if err != nil {
		return err
	}

	return nil
}
