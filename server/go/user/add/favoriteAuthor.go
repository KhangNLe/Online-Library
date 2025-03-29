package user

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func addFavoriteAuthor(c *gin.Context, query *sqlx.Tx, authorKey string, libId int) error {
	resp, err := query.Query(`SELECT author_id FROM Favorite_Author WHERE library_id = ?`,
		libId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not do a select statment for favorite_author. Error: %s", err)
		return err
	}
	defer resp.Close()

	for resp.Next() {
		var author_id string
		err = resp.Scan(&author_id)
		if err != nil {
			ErrorRespone(c, `
                We could not perform this action at this moment. Please try again later.
                `, http.StatusInternalServerError)
			log.Printf("Could not scan author_id. Error: %s", err)
			return err
		}
		if author_id == authorKey {
			ErrorRespone(c, `
                The author is already in your Favorite Author session.
                `, http.StatusBadRequest)
			return errors.New("Author is already in lib")
		}
	}

	_, err = query.Exec(`INSERT INTO Favorite_Author(library_id, author_id)
                            VALUES (?, ?)`, libId, authorKey)
	if err != nil {
		ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Could not insert author into favorite_author. Error: %s", err)
		return err
	}
	return nil
}
