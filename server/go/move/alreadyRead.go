package move

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func addToAlreadyRead(c *gin.Context, query *sqlx.Tx,
	libID int, bookKey string, userID int) error {
	book, err := query.Query(`SELECT P.book_id FROM Read_Book as P JOIN User_library as U
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
		book_id := ""
		err = book.Scan(&book_id)

		if err != nil {
			ErrorRespone(c, `
                We could not perform this action at the moment. 
                Please try again later.
                `, http.StatusBadRequest)
			log.Printf("Error while try to scan for book_id. Error: %s", err)
			return err
		}

		if book_id == bookKey {
			ErrorRespone(c, `
            The book is already in your Read session.
            `, http.StatusBadRequest)
			log.Printf("Error, book is already exist")
			return errors.New("Book already exist")
		}
	}
	_, err = query.Exec(`INSERT INTO Read_Book (library_id, book_id) VALUES (?, ?)`,
		libID, bookKey)
	if err != nil {
		ErrorRespone(c, ``, http.StatusBadRequest)
		log.Printf("Could not add book into the Read_Book db. Error: %s", err)
		return err
	}

	return nil
}
