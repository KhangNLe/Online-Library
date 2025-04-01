package mybook

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func moveToFavorite(c *gin.Context, query *sqlx.Tx,
	bookKey string, libId int) error {
	_, err := query.Exec(`INSERT INTO Favorite_Book (book_id, library_id)
        VALUES (?, ?)`, bookKey, libId)
	if err != nil {
		ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Unable to insert book into favorite. Error: %s", err)
		return err
	}

	return nil
}
