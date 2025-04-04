package recomend

import (
	"book/search"
	"fmt"
	"html"
	"log"
	"sort"
	"strconv"
	"strings"
)

func getBook(genre string, books *[]string, books_id *[]int) error {
	sort.Ints(*books_id)
	results, err := search.SearchBook(genre)
	if err != nil {
		log.Printf("Could not search for books from genre. Error: %s",
			err)
		return err
	}

	*books = append(*books, `
        <div class="display">
        <div class="search-display">
    `)
	count := 0
	for _, result := range results {
		if count == 12 {
			break
		}

		if checkForExistedBook(result.BookKey, books_id) {
			continue
		}
		var img string
		if result.CoverPic == 0 {
			img = "https://upload.wikimedia.org/wikipedia/commons/1/14/No_Image_Available.jpg?20200913095930"
		} else {
			img = html.EscapeString("https://covers.openlibrary.org/b/id/" + strconv.Itoa(result.CoverPic) + "-M.jpg")
		}
		*books = append(*books, fmt.Sprintf(`
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
			img,
			result.BookKey,
			result.Author_name[0],
			result.Author_key[0],
			img,
			result.BookKey,
			result.Title))
		count++
	}
	*books = append(*books, `
        </div>
        </div>
        `)
	return nil
}

func checkForExistedBook(book_id string, books_id *[]int) bool {
	book_id = strings.Replace(book_id, "/works/OL", "", 1)
	book_id = strings.Replace(book_id, "W", "", 1)
	num, err := strconv.Atoi(book_id)
	if err != nil {
		log.Printf("Could not convert the book_id: %s to int", book_id)
		return false
	}

	start := 0
	end := len(*books_id) - 1
	for start <= end {
		mid := (start + end) / 2
		if (*books_id)[mid] == num {
			return true
		} else if (*books_id)[mid] < num {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	return false
}
