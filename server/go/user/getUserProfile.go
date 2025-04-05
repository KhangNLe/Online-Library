package user

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type User struct {
	Fname string
	Lname string
	Email string
}

func getUserProfile(db *sqlx.DB, userId int) (User, error) {
	var user User

	resp, err := db.Query(`SELECT fname, lname, email FROM User_Profile 
		WHERE user_id = ?`, userId)
	if err != nil {
		log.Printf("Could not get the user profile. Error: %s", err)
		return User{}, err
	}
	defer resp.Close()
	if !resp.Next() {
		log.Printf("Could not find a user with the profile. Error: %s", err)
		return User{}, err
	}

	err = resp.Scan(&user.Fname, &user.Lname, &user.Email)
	if err != nil {
		log.Printf("Could not scan for user infos. Error: %s", err)
		return User{}, err
	}
	return user, nil
}
