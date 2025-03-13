package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Booksearch struct {
	ResultAmount int       `json:"numFound"`
	Result       []Results `json:"docs"`
}

type Results struct {
	Author_key   []string `json:"author_key"`
	Author_name  []string `json:"author_name"`
	CoverPic     int      `json:"cover_i"`
	Publish_year int      `json:"first_publish_year"`
	BookKey      string   `json:"key"`
	Title        string   `json:"title"`
}

func SearchBook(text string) []Results {
	url := "https://openlibrary.org/search.json?q="
	text = strings.ReplaceAll(text, " ", "+")

	var search Booksearch

	resp, err := http.Get(url + text)
	if err != nil {
		log.Printf("Could not fetch the %s url, erro: %s", (url + text), err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		log.Printf("Could not decode the api, err: %s", err)
	}

	return search.Result[:15]
}

func DisplaySearch(text string, c *gin.Context) {
	bookResult := SearchBook(text)

	var bookDisplay []string

	for index, book := range bookResult {
		if index%4 == 0 && index != 0 {
			bookDisplay = append(bookDisplay,
				fmt.Sprint(`
                    <div class="clear"></div>
                    `),
			)
		}
		bookDisplay = append(bookDisplay,
			fmt.Sprintf(`
                <div>
                <img src="%s">
                </div>
                <div class="book-title">
                    <button hx-post="/books" hx-swap="innerHTML"
                            hx-trigger="click" hx-target=".contents" value=%s>
                    %s</button>
                </div>
                `,
				("https://covers.openlibrary.org/b/id/"+strconv.Itoa(book.CoverPic)+"-M.jpg"),
				book.BookKey,
				book.Title,
			),
		)
	}
	c.Header("Content-Type", "text/html")
	c.String(200, strings.Join(bookDisplay, "\n"))
}
