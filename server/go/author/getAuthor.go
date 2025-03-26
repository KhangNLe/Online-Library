package author

import (
	"book/author"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	c.Header("Content-Type", "text/html")
	var authorPage []string
	authorPage = append(authorPage, fmt.Sprintf(`
            <div class="bookpageLeft">
                <div class="bookImg">
                    <img src="%s">
                </div>
                <div class="bookAction">
                    <div class="btn-group" role="group">
                        <div class="dropdown">
                            <button class="btn btn-success dropdown-toggle"
                                    type="button" data-bs-toggle="dropdown"
                                    aria-expanded="false"
                            >Add to Favorite</button>
                            <ul class="dropdown-menu">
                                <li><a class="dropdown-item" href="#">Block Author</a></li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
            <div class="bookpageRight">
                <div class="bookTitle">
        <p style="font-size: 25px;">Author: %s</p>
                </div>
                <div class="dob">
                    <span>
                        <h3 style="font-size: 17px; display: inline;">Birth Date:</h3>
                        <p style="display: inline;">    %s</p>
                    </spane>
                </div>
                <div class="dod">
                    <span>
                        <h3 style="font-size: 17px; display: inline;">Death Date:</h3>
                        <p style="display: inline;">    %s</p>
                    </span>
                </div>  
                <div class="bookDescription">
                    <h3>Bio:</h3>
                    <p>%s</p>
                </div>
            </div>
    `, author.Photo,
		author.Name,
		author.Birth,
		author.Death,
		author.Bio))

	c.String(200, strings.Join(authorPage, ""))
}

func linksDisplay(authorPage *[]string, links []Link) {
	for _, link := range author.Link {

	}
}
