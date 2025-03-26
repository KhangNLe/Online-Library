package search

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
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

var books []Results

func SearchBook(text string) []Results {
	url := "https://openlibrary.org/search.json?q="
	text = strings.ReplaceAll(text, " ", "+")
	fields := "&fields=author_key,author_name,cover_i,title,key"
	limit := "&limit=250"
	return getSearchBook(url + text + fields + limit)
}

func getSearchBook(searchStr string) []Results {

	var search Booksearch

	resp, err := http.Get(searchStr)
	if err != nil {
		log.Printf("Could not fetch the %s url, erro: %s", (searchStr), err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		log.Printf("Could not decode the api, err: %s", err)
	}
	books = search.Result

	return search.Result
}

func appendBooks(start, end, totalBook int, bookDisplay *[]string) {

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

		*bookDisplay = append(*bookDisplay,
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
                        href="#"
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
}

func appendSubject(bookDisplay *[]string) {
	*bookDisplay = append(*bookDisplay, `
        <nav class="navbar bg-body-tertiary">
            <div class="container-fluid" style="max-width: fit-content; margin-left: auto; margin-right:auto;">
                <form class="d-flex" role="search"
                        hx-post="/book-search"
                        hx-swap="innerHTML"
                        hx-target=".display"
                        hx-vals='{"page" : "1"}'
                        hx-push-url="true"
                        hx-on::after-request=" if (event.detail.xhr.status >= 400) { document.querySelector('.search-display').innerHTML = event.detail.xhr.responseText; }" 
                        >
            <input class="form-control me-2" size="75%" name="query" type="search" autocomplete="off" placeholder="Enter title or author" aria-label="Search">
                    <button id="searchBtn" class="btn btn-outline-success" type="submit">Search</button>
                </form>
            </div>
        </nav>
        <div class="display">
        `)
}
