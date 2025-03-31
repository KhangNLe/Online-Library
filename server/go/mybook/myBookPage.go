package mybook

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func MyBookPage(c *gin.Context, db *sqlx.DB, user string, option int) {
	userId, err := strconv.Atoi(user)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not convert user_id to int. Error: %s", err)
		return
	}

	query, err := db.Beginx()
	if err != nil {
		ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Could not start sqlx.Tx. Error: %s", err)
		return
	}
	defer func() {
		if err != nil {
			query.Rollback()
		}
	}()

	resp, err := query.Query(`SELECT library_id FROM User_library WHERE
                user_id = ?`, userId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Error when getting lib_id. Error: %s", err)
		return
	}
	defer resp.Close()

	if !resp.Next() {
		ErrorRespone(c, `
            Could not find library for this current user. Please contact the dev for this fix.
            `, http.StatusBadRequest)
		log.Printf("No library for the user_id %d", userId)
		return
	}

	var libId int
	err = resp.Scan(&libId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment. Please try again later.
            `, http.StatusBadRequest)
		log.Printf("Could not scan for lib_id. Error: %s", err)
		return
	}

	var myPage []string

	switch option {
	case 1:
		err = currentlyReading(c, query, libId, &myPage)
	}

	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Something wrong while trying to display book. Error: %s", err)
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, strings.Join(myPage, "\n"))
}

func ErrorRespone(c *gin.Context, msg string, status int) {
	c.Header("Content-Type", "text/html")
	c.String(status, fmt.Sprintf(`
        <br><p style="color: red; font-size: 15px;">
        %s
        </p></br>
        `, msg))
}
