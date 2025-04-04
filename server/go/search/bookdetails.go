package search

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Link struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

type SenDescription struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type BookDescription struct {
	Type  string
	Value string
}

func (bd *BookDescription) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		*bd = BookDescription{Value: str}
		return nil
	}

	type description struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	var obj description
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	*bd = BookDescription{
		Type:  obj.Type,
		Value: obj.Value,
	}
	return nil
}

type Book struct {
	Description  BookDescription `json:"description"`
	Title        string          `json:"title"`
	First_sen    SenDescription  `json:"first_sentence"`
	Subj_people  []string        `json:"subject_people"`
	Subjects     []string        `json:"subjects"`
	Links        Link            `json:"link"`
	Publish_year string          `json:"first_publish_date"`
	Key          string
	Cover        string
	Author       string
	AuthorKey    string
}

func BookDetail(c *gin.Context, db *sqlx.DB) {

	book := make(map[string]string)
	if err := c.Bind(&book); err != nil {
		c.Header("Content-Type", "text/html")
		c.Status(http.StatusInternalServerError)
		return
	}
	detail, ok := book["work"]
	if !ok {
		c.Header("Content-Type", "text/html")
		c.Status(http.StatusInternalServerError)
		return
	}

	author, ok := book["author"]
	if !ok {
		author = "Unknown author"
	}

	authorKey, ok := book["author_key"]
	if !ok {
		authorKey = ""
	}

	book_cover, ok := book["cover"]
	if !ok {
		book_cover = ""
	}

	reps, err := db.Query("SELECT * FROM Book_Detail WHERE book_id=?", detail)
	if err != nil {
		c.Header("Content-Type", "text/html")
		c.Status(http.StatusInternalServerError)
		return
	}
	defer reps.Close()

	if reps.Next() {
		var book Book
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
		genres, err := getBookGenre(db, book_id)
		if err != nil {
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
		book.AuthorKey = authorKey

		PrintBookDetail(book, c)
	} else {
		bookDetail, err := getBookDetail(detail)
		if err != nil {
			c.Header("Content-Type", "text/html")
			c.Status(http.StatusInternalServerError)
			return
		}

		bookDetail.Cover = book_cover
		bookDetail.Author = author
		bookDetail.AuthorKey = authorKey
		descriptions := strings.Split(bookDetail.Description.Value, "----------")
		noSource := strings.Split(descriptions[0], "([source]")
		bookDetail.Description.Value = noSource[0]
		err = addBookToDB(db, bookDetail)
		if err != nil {
			c.Header("Content-Type", "text/html")
			c.Status(http.StatusInternalServerError)
			return
		}
		PrintBookDetail(bookDetail, c)
	}
}

func getBookDetail(work string) (Book, error) {

	var book Book
	url := "https://openlibrary.org" + work + ".json"

	resp, err := http.Get(url)
	if err != nil {
		return Book{}, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&book)
	if err != nil {
		return Book{}, err
	}
	book.Key = work
	return book, nil
}

func addBookToDB(db *sqlx.DB, book Book) error {
	_, err := db.Exec(`INSERT INTO Book_Detail(book_id, cover_img,
                        title, description, first_setence, year_publish) VALUES
                        (?, ?, ?, ?, ?, ?)`,
		book.Key, book.Cover, book.Title, html.EscapeString(book.Description.Value),
		book.First_sen.Value, book.Publish_year)
	if err != nil {
		return err
	}

	for _, genre := range book.Subjects {
		var genre_id int
		err := db.Get(&genre_id, `SELECT genre_id FROM Genre
                                    WHERE genre_name=?`, genre)
		if err != nil {
			_, err = db.Exec(`INSERT INTO Genre(genre_name) VALUES (?)`, genre)
			if err != nil {
				return err
			}
		}
		_, err = db.Exec(`INSERT INTO Book_Genre(book_id, genre_id)
            SELECT ?, genre_id FROM Genre WHERE genre_name=?`, book.Key, genre)
		if err != nil {
			return err
		}
	}

	author := "/authors/" + book.AuthorKey
	_, err = db.Exec(`INSERT INTO Book(book_id, author_id) VALUES (?, ?)`,
		book.Key, author)
	if err != nil {
		return err
	}

	return nil
}

func getBookGenre(db *sqlx.DB, book_id string) ([]string, error) {
	var genres []string

	err := db.Select(&genres, `SELECT G.genre_name
                                FROM Genre AS G JOIN Book_Genre AS B
                                ON G.genre_id = B.genre_id
                                WHERE B.book_id=?`, book_id)
	if err != nil {
		return nil, err
	}

	return genres, nil
}

func presentingGenre(book Book) []string {
	var genres []string

	for idx, genre := range book.Subjects {
		if idx == 20 {
			break
		}
		genres = append(genres, fmt.Sprintf(`
            <a href="#"
                style="text-decoration: none; color: black;"
                hx-get="/book-search"
                hx-target=".contents"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-vals='{
                    "text":     "%s",
                    "page":     "1",
                    "subject":  "yes"
                        }'
                ><span class="genre"><span>%s</span></span></a>
        `, html.EscapeString(genre), genre))
	}

	return genres
}

func PrintBookDetail(bookDetail Book, c *gin.Context) {
	var details []string
	c.Header("Content-Type", "text/html")
	details = append(details, fmt.Sprintf(`
        <div class="contentContainer">
            <div class="contentLeft">
                <div class="bookImg">
                    <img src="%s">
                </div>
                <div class="bookAction">
                    <div class="btn-group" role="group"
                        style="max-height: 60px; margin-left: -12%% ;">
                        <button type="button" class="btn btn-success"
                            style="width: 250px;">
                            <a id="firstBookButton"
                            hx-get="/wantToRead/add"
                            hx-target=".responeMessage"
                            hx-swap="innerHTML"
                            hx-vals='{
                                "key": "%s"
                                }'
                            hx-on::after-request="
                                if (event.detail.xhr.status >= 400){
                                    document.querySelector('.responeMessage').innerHTML = event.detail.xhr.responseText;
                                }"
                            >Want to Read</a>
                        </button>
                        <div class="dropdown bookBtn btn-group"
                            style="width: 30px;">
                            <button class="btn btn-success dropdown-toggle"
                                    type="button" id="wantToRead" data-bs-toggle="dropdown"
                                    aria-expanded="false"
                            >
                            </button>
                            <ul class="dropdown-menu">
                                <li><a class="dropdown-item"
                                    href="#"
                                    hx-get="/alreadyRead/add"
                                    hx-target=".responeMessage"
                                    hx-swap="innerHTML"
                                    hx-trigger="click"
                                    hx-vals='{
                                        "key" : "%s"
                                        }'
                                    hx-on::after-request="
                                    if (event.detail.xhr.status >= 400){
                                        document.querySelector('.responeMessage').innerHTML = event.detail.xhr.responseText;
                                    }"
                                    >Add to Library</a></li>
                                <li><a class="dropdown-item" 
                                    href="#"
                                    hx-get="/reading/add"
                                    hx-target=".responeMessage"
                                    hx-swap="innerHTML"
                                    hx-trigger="click"
                                    hx-vals='{
                                        "key" : "%s"
                                        }'
                                    hx-on::after-request="
                                        if (event.detail.xhr.status >= 400){
                                        document.querySelector('.responeMessage').innerHTML = event.detail.xhr.responseText;
                                        }"
                                    >Reading</a></li>
                                <li><a class="dropdown-item" 
                                    href="#"
                                    hx-get="/favorite/add"
                                    hx-target=".responeMessage"
                                    hx-swap="innerHTML"
                                    hx-trigger="click"
                                    hx-vals='{
                                        "key" : "%s"
                                        }'
                                    hx-on::after-request="
                                        if (event.detail.xhr.status >= 400){
                                        document.querySelector('.responeMessage').innerHTML = event.detail.xhr.responseText;
                                        }"
                                    >Favorite</a></li>
                            </ul>
                    </div>
                </div>
                <div class="responeMessage"></div>
                </div>
            </div>
            <div class="contentRight">
                <div class="bookTitle">
                    <h3>%s</h3>
                    <p>Author: <a
                        id="author"
                        href="#"
                        hx-post="/author"
                        hx-target=".contentContainer"
                        hx-swap="innerHTML"
                        hx-push-url="/author/%s"
                        hx-vals='{
                            "key"           : "/authors/%s",
                            "bookKey"       : "%s",
                            "authorName"    : "%s"
                                }'
                        >
                        %s</a></p>
                </div>
                <div class="bookDescription">
                    <p>%s</p>
                </div>
                <div class="bookGenre">
                    <spane>Genres:</span>
                    <ul class="genreList">
        `, bookDetail.Cover,
		bookDetail.Key,
		bookDetail.Key,
		bookDetail.Key,
		bookDetail.Key,
		bookDetail.Title,
		bookDetail.AuthorKey,
		bookDetail.AuthorKey,
		bookDetail.Key,
		bookDetail.Author,
		bookDetail.Author,
		html.EscapeString(bookDetail.Description.Value),
	))

	genres := presentingGenre(bookDetail)
	genreList := strings.Join(genres, "\n")
	details = append(details, genreList)
	details = append(details, `
                    </ul>
                </div>
            </div>
        </div>
        `)
	c.String(200, strings.Join(details, ""))
}
