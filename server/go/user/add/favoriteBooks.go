package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func addToFavoriteBook(c *gin.Context, query *sqlx.Tx, libId int, bookKey string) error {
	resp, err := query.Query(`SELECT book_id FROM Favorite_Book 
                WHERE library_id = ?`, libId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not perform Select on Favorite_Book. Error: %s", err)
		return err
	}
	defer resp.Close()

	for resp.Next() {
		var book_id string
		err = resp.Scan(&book_id)
		if err != nil {
			ErrorRespone(c,
				`We could not perform this action at the moment. Pleasy try again.
                `, http.StatusInternalServerError)
			log.Printf("Could not scan book_id. Error: %s", err)
			return err
		}

		if book_id == bookKey {
			ErrorRespone(c, `
                The book is already in your Favorite Book.
                `, http.StatusBadRequest)
			return errors.New("Book is already in the Favorite book.")
		}
	}
	return nil
}
