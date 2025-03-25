package search

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func LoadingBookDetail(c *gin.Context, bookKey string, db *sqlx.DB) {
	var reloadPage []string
	bookKey = "/works/" + bookKey

	genres, err := getBookGenre(db, bookKey)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Something wrong with the db while trying to get genres. Error: %s", err)
		return
	}
	var book Book

	reps, err := db.Query(`SELECT A.name FROM Author AS A JOIN Book as B
                            ON A.author_id = B.author_id
                            WHERE B.book_id = ?`, bookKey)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Something wrong while getting author name. Error: %s", err)
		return
	}

	author := ""
	if reps.Next() {
		err = reps.Scan(&author)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Printf("Could not find an author name for the book key. Error: %s", err)
			return
		}
	}

	reps, err = db.Query("SELECT * FROM Book_Detail WHERE book_id=?", bookKey)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.Status(http.StatusInternalServerError)
		return
	}
	defer reps.Close()

	if reps.Next() {
		book_id := ""
		cover_img := ""
		title := ""
		description := ""
		first_sen := ""
		year_publish := ""
		err = reps.Scan(&book_id, &cover_img, &title,
			&description, &first_sen, &year_publish)
		if err != nil {
			c.Header("Content-Type", "text/html")
			c.Status(http.StatusInternalServerError)
			return
		}

		book.Cover = cover_img
		book.Key = book_id
		book.Title = title
		book.Description.Value = description
		book.First_sen.Value = first_sen
		book.Publish_year = year_publish
		book.Subjects = genres
		book.Author = author
	}

	reloadBook(&reloadPage, book)

	c.Header("Content-Type", "text/html")
	c.String(200, strings.Join(reloadPage, "\n"))
}

func reloadBook(page *[]string, book Book) {
	*page = append(*page, fmt.Sprintf(`
            <div class="bookpageLeft">
                <div class="bookImg">
                    <img src="%s">
                </div>
                <div class="bookAction">
                    <div class="btn-group" role="group">
                        <button type="button" class="btn btn-success">Want to Read</button>
                        <div class="dropdown">
                            <button class="btn btn-success dropdown-toggle"
                                    type="button" data-bs-toggle="dropdown"
                                    aria-expanded="false"
                            ></button>
                            <ul class="dropdown-menu">
                                <li><a class="dropdown-item" href="#">Add to Library</a></li>
                                <li><a class="dropdown-item" href="#">Reading</a></li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
            <div class="bookpageRight">
                <div class="bookTitle">
                    <h3>%s</h3>
                    <p>Author: %s</p>
                </div>
                <div class="bookDescription">
                    <p>%s</p>
                </div>
                <div class="bookGenre">
                    <spane>Genres:</span>
                    <ul class="genreList">
        `,
		book.Cover,
		book.Title,
		book.Author,
		html.EscapeString(book.Description.Value),
	))

	genres := presentingGenre(book)
	*page = append(*page, strings.Join(genres, ""))
	*page = append(*page, `
                    </ul>
                </div>
            </div>
    `)
}
