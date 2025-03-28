package search

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Display struct {
	Author_key  string
	Author_name string
	Title       string
	IMG         int
	BookKey     string
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
		log.Printf("Error with getting hx-vals. Error: %s", err)
	}

	currPage, _ := page["page"]
	pageNum, _ := strconv.Atoi(currPage)
	start := 21 * (pageNum - 1)
	end := 21 * pageNum

	var text string
	if result, ok := page["text"]; !ok {
		text = c.PostForm("query")
	} else {
		text = result
	}

	var totalBook int

	books, err := SearchBook(text)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.String(200, fmt.Sprintln(`
            <p>We could not fetch your books at the current moment, please try again in a bit.</p>
            `))
		log.Printf("Error while searching book. Error %s", err)
		return
	}
	totalBook = len(books)

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

	_, ok := page["subject"]
	if ok {
		appendSubject(&bookDisplay)
	}
	bookDisplay = append(bookDisplay, `<div class="search-display">`)

	appendBooks(start, end, totalBook, &bookDisplay)

	addingPageBtn(&bookDisplay, pageNum, totalPage, text)
	if ok {
		bookDisplay = append(bookDisplay, "</div>")
	}

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
	if totalPage < 4 {
		for i := 1; i < totalPage; i++ {
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
		if pageNum < 3 {
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
		} else if pageNum+1 >= totalPage {
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
	}

	*bookDisplay = append(*bookDisplay, `
        </div>
        </div>
        </div>
        `)
}
