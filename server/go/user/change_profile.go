package user

import (
	"book/mybook"
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
			We could not make this change at the this moment.
			Please contact the dev for the problem.
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
			We could not make this change at the this moment.
			Please contact the dev for the problem.
			`, http.StatusInternalServerError)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	ok := checkProperEmail(email)
	if !ok {
		log.Printf("Inproper email address of %s", email)
		mybook.ErrorRespone(c, `
			Incorrect email format, please try again.
			`, http.StatusBadRequest)
		return
	}

	userExist, err := checkForExistingUser(tx, userId)
	if err != nil {
		log.Printf("Could not find existed user. Error: %s", err)
		mybook.ErrorRespone(c, `
			We could not make this change at the this moment.
			Please contact the dev for the problem.
			`, http.StatusInternalServerError)
		return
	}

	if !userExist {
		err = setNewProfile(tx, fname, lname, email, userId)
	} else {
		err = updateProfile(tx, fname, lname, email, userId)
	}
	if err != nil {
		log.Printf("Could not update profile. Error: %s", err)
		mybook.ErrorRespone(c, `
			We could not make this change at the this moment.
			Please contact the dev for the problem.
				`, http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Could not commit into db. Error: %s", err)
		mybook.ErrorRespone(c, `
			We could not make this change at the this moment.
			Please contact the dev for the problem.
			`, http.StatusInternalServerError)
		return
	}

	UserProfile(c, db, user)
}

func updateProfile(tx *sqlx.Tx, fname, lname, email string,
	userId int) error {

	_, err := tx.Exec(`UPDATE User_Profile 
			SET fname = ?, lname = ?, email = ? 
			WHERE user_id = ?`, fname, lname, email, userId)
	if err != nil {
		log.Printf("could not update user info. Error %s", err)
		return err
	}

	return nil
}

func setNewProfile(tx *sqlx.Tx, fname, lname, email string,
	userId int) error {

	_, err := tx.Exec(`INSERT INTO User_Profile 
		VALUES (?, ?, ?, ?)`, userId, fname, lname, email)
	if err != nil {
		log.Printf("Could not insert in new user_profile. Error: %s", err)
		return err
	}
	return nil
}

func checkForExistingUser(tx *sqlx.Tx, userId int) (bool, error) {
	resp, err := tx.Query(`SELECT * FROM User_Profile 
			WHERE user_id = ?`, userId)
	if err != nil {
		log.Printf("Could not find User_profile with userid. Error: %s", err)
		return false, err
	}
	defer resp.Close()
	if !resp.Next() {
		return false, nil
	}
	return true, nil
}

func checkProperEmail(email string) bool {
	emailRegex := `(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`

	regex := regexp.MustCompile(emailRegex)

	return regex.MatchString(email)
}
