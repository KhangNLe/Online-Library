package mybook

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func currentlyReading(c *gin.Context, query *sqlx.Tx,
	libId int, myBooksPage *[]string) error {
	resp, err := query.Query(`SELECT B.cover_img, B.title FROM Book_Detail as B JOIN Reading as R
                    ON B.book_id = R.book_id
                    WHERE R.library_id = ?`, libId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action that this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not get books from Book_detail. Error: %s", err)
		return err
	}
	defer resp.Close()

	for resp.Next() {
		img := ""
		title := ""

		err = resp.Scan(&img, &title)
		if err != nil {
			ErrorRespone(c, ``, http.StatusInternalServerError)
			log.Printf("Could not get img and book title. Error: %s", err)
		}

		*myBooksPage = append(*myBooksPage, fmt.Sprintf(`

        `))
	}

	return nil
}
