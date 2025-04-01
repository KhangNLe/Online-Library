package mybook

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func moveToToRead(c *gin.Context, query *sqlx.Tx,
	from, bookKey string, libId int) error {
	var err error
	switch from {
	case "reading":
		_, err = query.Exec(`DELETE FROM Reading 
            WHERE book_id = ? AND library_id = ?`, bookKey, libId)
	case "finish":
		_, err = query.Exec(`DELETE FROM Read_Book 
            WHERE book_id = ? AND library_id = ?`, bookKey, libId)
	default:
		err = errors.New("Did not match any of the options")
	}
	if err != nil {
		ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Could not delete book from previous container %s. Error: %s",
			from, err)
	}

	_, err = query.Exec(`INSERT INTO Planning_to_Read (library_id, book_id)
        VALUES (?, ?)`, libId, bookKey)
	if err != nil {
		ErrorRespone(c, ``, http.StatusBadRequest)
		log.Printf("Could not insert into planning to read. Error: %s", err)
	}
	return nil
}
