package move

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func addToReading(c *gin.Context, query *sqlx.Tx, libID int,
	bookKey string, userID int) error {
	book, err := query.Query(`SELECT P.book_id FROM Reading as P JOIN User_library as U
                        ON P.library_id = U.library_id
                        WHERE U.user_id = ?`, userID)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment, please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Error: %s", err)
		return err
	}
	defer book.Close()

	for book.Next() {
		bookId := ""
		err = book.Scan(&bookId)
		if err != nil {
			ErrorRespone(c, ``, http.StatusInternalServerError)
			log.Printf("Could not scan book_id. Error: %s", err)
			return err
		}

		if bookId == bookKey {
			ErrorRespone(c, `
            The book is already in your Reading session".
            `, http.StatusBadRequest)
			log.Printf("Error, book is already exist")
			return errors.New("There is already a book there")
		}
	}
	_, err = query.Exec(`INSERT INTO Reading (library_id, book_id) VALUES (?, ?)`,
		libID, bookKey)
	if err != nil {
		ErrorRespone(c, ``, http.StatusBadRequest)
		log.Printf("Could not adding book into Reading, Error: %s", err)
		return err
	}
	return nil
}
