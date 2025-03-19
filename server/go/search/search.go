package search

import (
	"encoding/json"
	"fmt"
	"html"
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
	Author_key  []string `json:"author_key"`
	Author_name []string `json:"author_name"`
	CoverPic    int      `json:"cover_i"`
	BookKey     string   `json:"key"`
	Title       string   `json:"title"`
}

func SearchBook(text string) []Results {
	url := "https://openlibrary.org/search.json?q="
	text = strings.ReplaceAll(text, " ", "+")
	fields := "&fields=author_key,author_name,cover_i,title,key"

	var search Booksearch

	resp, err := http.Get(url + text + fields)
	if err != nil {
		log.Printf("Could not fetch the %s url, erro: %s", (url + text), err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		log.Printf("Could not decode the api, err: %s", err)
	}

	return search.Result
}

func DisplaySearch(c *gin.Context) {
	text := c.PostForm("query")
	bookResult := SearchBook(text)

	if len(bookResult) == 0 {
		c.Header("Context-Type", "text/html")
		c.String(200, fmt.Sprintf(`
                <p style="font-size:24px;">We could not find any book with the search "%s". Please try another input
                </p>
                `, text))
		return
	}

	var bookDisplay []string
	count := 0
	for _, book := range bookResult {
		if count == 25 {
			break
		}
		if count%4 == 0 {
			bookDisplay = append(bookDisplay,
				fmt.Sprint(`
                    <div class="clear"></div>
                    `),
			)
		}
		if len(book.Author_name) == 0 {
			continue
		}
		var bookPic string
		if book.CoverPic == 0 {
			bookPic = "https://upload.wikimedia.org/wikipedia/commons/1/14/No_Image_Available.jpg?20200913095930"
		} else {
			bookPic = html.EscapeString("https://covers.openlibrary.org/b/id/" + strconv.Itoa(book.CoverPic) + "-M.jpg")
		}

		count++
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
				html.EscapeString(book.BookKey),
				book.Author_name[0],
				book.Author_key[0],
				bookPic,
				book.BookKey,
				book.Title,
			),
		)
	}
	c.Header("Content-Type", "text/html")
	c.String(200, strings.Join(bookDisplay, "\n"))
}
