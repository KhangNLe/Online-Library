package user

import (
	password "book/loginsignup"
	"book/mybook"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func ChangePass(c *gin.Context, db *sqlx.DB, user string) error {
	userId, err := strconv.Atoi(user)
	if err != nil {
		log.Printf("Could not convert %s into integer. Error: %s", user, err)
		mybook.ErrorRespone(c, `
			Could not change your password at the moment. Please contact the dev for the problem.
			`, http.StatusInternalServerError)
		return err
	}

	pass := c.PostForm("currPass")
	newPass := c.PostForm("newPass")

	hashPass, err := getHashPass(db, userId)
	if err != nil {
		log.Printf("Could not get pass_hash. Error: %s", err)
		mybook.ErrorRespone(c, `
			Could not change your password at the moment. Please contact the dev for the problem.
			`, http.StatusInternalServerError)
		return err
	}

	if !password.ComparePass(hashPass, pass) {
		log.Printf("Incorrect pass_hash")
		mybook.ErrorRespone(c, `
			The password you entered for curernt password is incorrect.
			`, http.StatusBadRequest)
		return err
	}

	pass, err = password.HashPass(newPass)
	if err != nil {
		log.Printf("Could not generate the new password. Error: %s", err)
		mybook.ErrorRespone(c, `
			We could not change your password at the moment. Please contact the dev for the problem.
			`, http.StatusInternalServerError)
		return err
	}

	err = updatePassword(db, pass, userId)
	if err != nil {
		mybook.ErrorRespone(c, `
			We could not change your password at the moment. Please contact the dev for this problem.
			`, http.StatusInternalServerError)
		return err
	}
	return nil
}

func getHashPass(db *sqlx.DB, userId int) (string, error) {
	resp, err := db.Query(`SELECT pass_hash FROM User 
		WHERE user_id = ?`, userId)
	if err != nil {
		log.Printf("Could not get the hash pass from user. Error: %s", err)
		return "", err
	}
	defer resp.Close()
	var hashPass string

	if !resp.Next() {
		log.Printf("Could not find a pass_hash for user_id: %d", userId)
		return "", errors.New("Could not find hash_pass")
	}

	err = resp.Scan(&hashPass)
	if err != nil {
		log.Printf("Could not scan for pass. Error: %s", err)
		return "", err
	}
	return hashPass, nil
}

func updatePassword(db *sqlx.DB, newPass string, userId int) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Printf("Could not get sqlx.Tx. Error: %s", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`UPDATE User SET pass_hash = ? 
		WHERE user_id = ?`, newPass, userId)
	if err != nil {
		log.Printf("Could not update new password for User. Error: %s", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Could not commit the change. Error: %s", err)
		return err
	}
	return nil
}
