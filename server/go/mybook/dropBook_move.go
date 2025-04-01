package mybook

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func dropBook(c *gin.Context, query *sqlx.Tx,
	from, bookKey string, libId int) error {
	var err error
	switch from {
	case "reading":
		_, err = query.Exec(`DELETE FROM Reading 
                WHERE library_id = ? AND book_id = ?`, libId, bookKey)
	case "toread":
		_, err = query.Exec(`DELETE FROM Planning_to_Read 
                WHERE library_id = ? AND book_id = ?`, libId, bookKey)
	case "favorite":
		_, err = query.Exec(`DELETE FROM Favorite_Book 
            WHERE library_id = ? AND book_id = ?`, libId, bookKey)
	case "finish":
		_, err = query.Exec(`DELETE FROM Read_Book 
            WHERE library_id = ? AND book_id = ?`, libId, bookKey)
	default:
		err = errors.New(fmt.Sprintf("The from: %s did not match any option", from))
	}

	if err != nil {
		ErrorRespone(c, ``, http.StatusBadRequest)
		log.Printf("Could not drop book from %s. Error: %s", from, err)
		return err
	}
	return nil
}
