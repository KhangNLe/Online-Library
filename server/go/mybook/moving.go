package mybook

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func MovingBooks(c *gin.Context, db *sqlx.DB, dst string, user string) {
	userId, err := strconv.Atoi(user)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusBadRequest)
		log.Printf("User session id was not an int. Got: %s. Error: %s",
			user, err)
	}

	bookVals := make(map[string]string)
	err = c.Bind(&bookVals)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusBadRequest)
		log.Printf("Could not get the value from action. Error: %s", err)
		return
	}

	from, ok := bookVals["from"]
	if !ok {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Println("Could not find 'from' in hx-vals")
		return
	}

	bookKey, ok := bookVals["key"]
	if !ok {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Println("Could not find book_key in hx-vals")
		return
	}

	query, err := db.Beginx()
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Error while transfer it to sqlx.Tx. Error: %s", err)
		return
	}
	defer func() {
		if err != nil {
			query.Rollback()
		}
	}()

	resp, err := query.Query(`SELECT library_id FROM User_library
            WHERE user_id = ?`, userId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusBadRequest)
		log.Printf("Error requesting library_id. Error: %s", err)
		return
	}
	defer resp.Close()

	var lib_id int
	if !resp.Next() {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusBadRequest)
		log.Printf("Could not find library_id for user_id: %s", user)
		return
	}

	err = resp.Scan(&lib_id)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not scan for library_id. Error: %s", err)
		return
	}

	switch dst {
	case "finish":
		err = moveToFinishReading(c, query, from, bookKey, lib_id)
	case "favorite":
		err = moveToFavorite(c, query, bookKey, lib_id)
	case "toread":
		err = moveToToRead(c, query, from, bookKey, lib_id)
	case "drop":
		err = dropBook(c, query, from, bookKey, lib_id)
	default:
		err = errors.New("Destination did not match any option")
	}

	if err != nil {
		log.Printf("Error happened. Error: %s", err)
		return
	}

	err = query.Commit()
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not commit the database change. Error: %s", err)
		return
	}
}
