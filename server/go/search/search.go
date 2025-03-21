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
	limit := "&limit=250"

	var search Booksearch

	resp, err := http.Get(url + text + fields + limit)
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
	page := make(map[string]string)
	err := c.Bind(&page)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusInternalServerError, `
            <p style="font-size:24px;">We are experiencing some techincal difficulty, please try again later
            </p>
            `)
	}
	currPage, _ := page["page"]
	pageNum, _ := strconv.Atoi(currPage)
	start := 21 * (pageNum - 1)
	end := 21 * pageNum
	var text string
	if pageNum == 1 {
		text = c.PostForm("query")
	} else {
		text, _ = page["text"]
	}
	books := SearchBook(text)

	totalBook := len(books)
	totalPage := 0
	if totalBook%21 == 0 {
		totalPage = totalBook / 21
	} else {
		totalPage = (totalBook / 21) + 1
	}

	if totalBook == 0 {
		c.Header("Context-Type", "text/html")
		c.String(200, fmt.Sprintf(`
                <p style="font-size:24px;">We could not find any book with the search "%s". Please try another input
                </p>
                `, text))
		return
	}

	var bookDisplay []string

	bookDisplay = append(bookDisplay, `<div class="search-display">`)

	for i := start; i < end; i++ {
		if i == totalBook {
			break
		}

		if len(books[i].Author_name) == 0 {
			continue
		}
		var bookPic string
		if books[i].CoverPic == 0 {
			bookPic = "https://upload.wikimedia.org/wikipedia/commons/1/14/No_Image_Available.jpg?20200913095930"
		} else {
			bookPic = html.EscapeString("https://covers.openlibrary.org/b/id/" + strconv.Itoa(books[i].CoverPic) + "-M.jpg")
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
                        hx-replace-url="/book%s"
                        hx-push-url="true"
                    >
                    %s</a>
                </div>`,
				bookPic,
				html.EscapeString(books[i].BookKey),
				books[i].Author_name[0],
				books[i].Author_key[0],
				bookPic,
				books[i].BookKey,
				books[i].Title,
			),
		)
	}

	addingPageBtn(&bookDisplay, pageNum, totalPage, text)

	c.Header("Content-Type", "text/html")
	c.String(200, strings.Join(bookDisplay, "\n"))
}

func addingPageBtn(bookDisplay *[]string, pageNum int, totalPage int, text string) {
	*bookDisplay = append(*bookDisplay, "</div>")
	*bookDisplay = append(*bookDisplay, `
        <div class="pageBtns">
        <div class="btn-toolbar" role="toolbar">
        <div class="btn-group me-2" role="group">
        `)
	if pageNum < 4 {
		for i := 1; i < 4; i++ {
			*bookDisplay = append(*bookDisplay, fmt.Sprintf(`
                <button type="button" 
                class="pageBtn btn btn-primary"
                hx-get="/book-search"
                hx-target=".display"
                hx-swap="innerHTML"
                hx-vals='{"page": "%d",
                            "text": "%s"
                        }'
                hx-swap-url="/book-search/page/%d"
                hx-push-url="true"
                >%d</button>
                `, i, text, i, i))
		}
		*bookDisplay = append(*bookDisplay, `
            <button type="click" class="pageBtn btn btn-primary">...</button>
            `)
		*bookDisplay = append(*bookDisplay, fmt.Sprintf(`
            <button type="button" 
            class="pageBtn btn btn-primary"
            hx-get="/book-search"
            hx-target=".display"
            hx-swap="innerHTML"
            hx-vals='{"page": "%d",
                        "text": "%s"
                    }'
            hx-swap-url="/book-search/page/%d"
            hx-push-url="true"
            >%d</button>
            `, totalPage, text, totalPage, totalPage))
	} else if pageNum+3 >= totalPage {
		*bookDisplay = append(*bookDisplay, fmt.Sprintf(`
            <button type="button" 
            class="pageBtn btn btn-primary"
            hx-get="/book-search"
            hx-target=".display"
            hx-swap="innerHTML"
            hx-vals='{"page": "1",
                        "text": "%s"
                    }'
            hx-swap-url="/book-search/page/1"
            hx-push-url="true"
            >1</button>
            `, text))
		*bookDisplay = append(*bookDisplay, `
            <button type="click" class="pageBtn btn btn-primary">...</button>
            `)
		for i := totalPage - 2; i <= totalPage; i++ {
			*bookDisplay = append(*bookDisplay, fmt.Sprintf(`
                <button type="button" 
                class="pageBtn btn btn-primary"
                hx-get="/book-search"
                hx-target=".display"
                hx-swap="innerHTML"
                hx-vals='{"page": "%d",
                            "text": "%s"
                        }'
                hx-swap-url="/book-search/page/%d"
                hx-push-url="true"
                >%d</button>
                `, i, text, i, i))
		}

	} else {
		*bookDisplay = append(*bookDisplay, fmt.Sprintf(`
            <button type="button" 
            class="pageBtn btn btn-primary"
            hx-get="/book-search"
            hx-target=".display"
            hx-swap="outterHTML"
            hx-vals='{"page": "1",
                        "text": "%s"
                    }'
            hx-swap-url="/book-search/page/1"
            hx-push-url="true"
            >1</button>
            `, text))
		*bookDisplay = append(*bookDisplay, `
            <button type="click" class="pageBtn btn btn-primary">...</button>
            `)
		for i := pageNum - 1; i <= pageNum+1; i++ {
			*bookDisplay = append(*bookDisplay, fmt.Sprintf(`
                <button type="button" 
                class="pageBtn btn btn-primary"
                hx-get="/book-search"
                hx-target=".display"
                hx-swap="innerHTML"
                hx-vals='{"page": "%d",
                            "text": "%s"
                        }'
                hx-swap-url="/book-search/page/%d"
                hx-push-url="true"
                >%d</button>
                `, i, text, i, i))

		}

		*bookDisplay = append(*bookDisplay, `
            <button type="click" class="pageBtn btn btn-primary">...</button>
            `)
		*bookDisplay = append(*bookDisplay, fmt.Sprintf(`
            <button type="button" 
            class="pageBtn btn btn-primary"
            hx-get="/book-search"
            hx-target=".display"
            hx-swap="innerHTML"
            hx-vals='{"page": "%d",
                        "text": "%s"
                    }'
            hx-swap-url="/book-search/page/%d"
            hx-push-url="true"
            >%d</button>
            `, totalPage, text, totalPage, totalPage))
	}

	*bookDisplay = append(*bookDisplay, `
        </div>
        </div>
        </div>
        `)
}
