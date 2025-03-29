package user

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func addToBlockAuthor(c *gin.Context, query *sqlx.Tx, libId int, authorKey string) error {
	resp, err := query.Query(`SELECT author_id FROM Block_Author 
                    WHERE library_id = ?`, libId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not perform a select query from block_author. Error: %s",
			err)
		return err
	}
	defer resp.Close()

	for resp.Next() {
		var author_id string
		err = resp.Scan(&author_id)
		if err != nil {
			ErrorRespone(c, `
            We could not perform this action at the moment. Please try again later.
                `, http.StatusInternalServerError)
			log.Printf("Could not scan for author_id. Error: %s", err)
			return err
		}

		if author_id == authorKey {
			ErrorRespone(c, `
                The author is already in the your block list.
                `, http.StatusBadRequest)
			log.Printf("Author is alread in the block lib")
			return errors.New("Athor is already in the block_library")
		}
	}

	_, err = query.Query(`INSERT INTO Block_Author (library_id, author_id) 
                VALUES (?, ?)`, libId, authorKey)
	if err != nil {
		ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Could not insert into block_author. Error: %s", err)
		return err
	}

	return nil
}
