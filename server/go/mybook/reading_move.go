package mybook

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func moveBookToReading(c *gin.Context, query *sqlx.Tx, libId int,
	bookKey, from string) error {

	var err error
	switch from {
	case "toread":
		_, err = query.Exec(`DELETE FROM Planning_to_Read 
            WHERE library_id = ? AND book_id = ?`, libId, bookKey)
	case "finish":
		_, err = query.Exec(`DELETE FROM Read_Book 
            WHERE library_id = ? AND book_id = ?`, libId, bookKey)
	case "favorites":
		break
	default:
		err = errors.New(fmt.Sprintf(`
            The %s did not match with any options.
        `, from))
	}

	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment.
            Please try again later
            `, http.StatusBadRequest)
		log.Printf("Error happened when tried to delete item. Error: %s", err)
		return err
	}

	_, err = query.Exec(`INSERT INTO Reading (book_id, library_id) 
            VALUES (?, ?)`, bookKey, libId)

	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment.
            Please try again later.
            `, http.StatusBadRequest)
		log.Printf("Could not insert into reading. Error: %s", err)
		return err
	}

	return nil
}
