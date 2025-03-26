package author

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetAuthor(c *gin.Context, db *sqlx.DB) {
	booksVals := make(map[string]string)

	err := c.Bind(&booksVals)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Could not get the info from bookDetail. Error: %s", err)
		return
	}
	authorKey, ok := booksVals["key"]
	if !ok {
		c.Status(http.StatusInternalServerError)
		log.Println("Could not find an author key")
		return
	}

	bookKey, ok := booksVals["bookKey"]
	if !ok {
		c.Status(http.StatusInternalServerError)
		log.Println("Unable to get the book key for this author")
		return
	}

	author, err := findAuthor(authorKey, bookKey, db)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Error while looking up author. Error: %s", err)
		return
	}

	printAuthor(c, author)
}

func printAuthor(c *gin.Context, author Author) {

}
