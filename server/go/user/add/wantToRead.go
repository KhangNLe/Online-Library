package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"strconv"
)

func WantToRead(userId string, c *gin.Context, db *sqlx.DB) {
	var hxVals map[string]string

	err := c.Bind(&hxVals)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("There are no hx-vals on button. Error: %s", err)
		return
	}

	bookKey, _ := hxVals["book"]

	query, err := db.Beginx()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Error starting db.Beginx. Error: %s", err)
		return
	}
	defer func() {
		if err != nil {
			query.Rollback()
		}
	}()

	num, _ := strconv.Atoi(userId)

	libraryId, err := query.Query(`SELECT library_id FROM User_library WHERE user_id = ?`, num)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Error with library_id search. Error: %s", err)
		return
	}
	defer libraryId.Close()

	var libID int

	if !libraryId.Next() {
		c.Status(http.StatusInternalServerError)
		log.Println("There is no library_id for this user")
		return
	}
	err = libraryId.Scan(&libID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Something wrong why trying to scan for library_id. Error: %s", err)
		return
	}

	book, err := query.Query(`SELECT P.book_id FROM Planning_to_Read as P JOIN User_library as U
                        ON P.library_id = U.library_id
                        WHERE U.user_id = ?`, num)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Error: %s", err)
		return
	}
	defer book.Close()

	if book.Next() {
		c.Status(http.StatusConflict)
		return
	}

	_, err = query.Exec(`INSERT INTO Planning_to_Read(library_id, book_id) VALUES (?, ?)`, libID, num)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Could not make the insert. Error : %s", err)
		return
	}

	err = query.Commit()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Could not commit change to db. Error: %s", err)
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
    <p>Book is in the Want to Read Library</p>
    `)

}
