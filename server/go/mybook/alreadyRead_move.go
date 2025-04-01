package mybook

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func moveToFinishReading(c *gin.Context, query *sqlx.Tx,
	from, bookKey string, libId int) error {

	var err error
	switch from {
	case "reading":
		_, err = query.Exec(`DELETE FROM Reading 
            WHERE book_id = ? AND library_id = ?`, bookKey, libId)
	case "toread":
		_, err = query.Exec(`DELETE FROM Planning_to_Read 
            WHERE book_id = ? AND library_id = ?`, bookKey, libId)
	default:
		err = errors.New("Could not find the proper from location")
	}

	if err != nil {
		ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Error occurd when trying to delete book. Error: %s", err)
		return err
	}

	_, err = query.Exec(`INSERT INTO Read_Book (library_id, book_id) 
        VALUES (?, ?)`, libId, bookKey)
	if err != nil {
		ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Could not move book into Read_Book. Error: %s", err)
		return err
	}
	return nil
}
