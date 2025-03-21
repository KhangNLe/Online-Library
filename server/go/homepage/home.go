package homepage

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Homepage(c *gin.Context) {
	books, err := TopFiction()
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusBadGateway, fmt.Sprintf(`
            <p>Could not load the top recommended book at this moment</p>
                `))
		return
	}

	var bookDisplay []string
	bookDisplay = append(bookDisplay, `<div class="search-display">`)
	log.Println(len(books))
	for _, book := range books {

		if len(book.AuthorName) == 0 {
			book.AuthorName = append(book.AuthorName, "Unknown Author")
		}
		if len(book.AuthorKey) == 0 {
			continue
		}
		var bookPic string
		if book.IMG == 0 {
			bookPic = "https://upload.wikimedia.org/wikipedia/commons/1/14/No_Image_Available.jpg?20200913095930"
		} else {
			bookPic = html.EscapeString("https://covers.openlibrary.org/b/id/" + strconv.Itoa(book.IMG) + "-M.jpg")
		}

		bookDisplay = append(bookDisplay,
			fmt.Sprintf(`
                <div class="books">
                <img src="%s">
                    <a hx-post="/book" hx-swap="innerHTML"
                        hx-trigger="click"
                        hx-target=".contents" 
                        hx-vals='{
                            "work":     "%s",
                            "author":   "%s",
                            "author_key":   "%s",
                            "cover":    "%s"}'
                        href="#"
                        hx-replace-url="/book%s"
                        hx-push-url="true"
                    >
                    %s</a>
                </div>`,
				bookPic,
				html.EscapeString(book.Key),
				book.AuthorName[0],
				book.AuthorKey[0],
				bookPic,
				book.Key,
				book.Title,
			),
		)
	}
	bookDisplay = append(bookDisplay, "</div>")

	c.Header("Content-Type", "text/html")
	c.String(200, strings.Join(bookDisplay, "\n"))
}
