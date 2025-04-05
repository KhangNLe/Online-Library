package move

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func wantToRead(c *gin.Context, query *sqlx.Tx, libID int,
	bookKey string, userID int) error {
	book, err := query.Query(`SELECT P.book_id FROM Planning_to_Read as P JOIN User_library as U
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
                Please try again later
                `, http.StatusInternalServerError)
			log.Printf("Error while scanning. Error: %s", err)
			return err
		}

		if book_id == bookKey {
			ErrorRespone(c, `
            The book is already in your Wanting to Read session.
            `, http.StatusBadRequest)
			log.Printf("Error, book is already exist")
			return errors.New("Book is already exist")
		}
	}
	_, err = query.Exec(`INSERT INTO Planning_to_Read(library_id, book_id) VALUES (?, ?)`,
		libID, bookKey)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment, please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not make the insert. Error : %s", err)
		return err
	}
	return nil
}
