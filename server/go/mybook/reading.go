package mybook

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func currentlyReading(c *gin.Context, query *sqlx.Tx,
	libId int, myBooksPage *[]string) error {
	loadOptions(myBooksPage)

	resp, err := query.Query(`SELECT Bd.title, Bd.book_id, 
                        Bd.cover_img, A.name, A.author_id
                        FROM Book B JOIN Author A ON B.author_id = A.author_id
                        JOIN Book_Detail Bd ON B.book_id = Bd.book_id
                        WHERE B.book_id
                        IN (SELECT book_id FROM Reading WHERE library_id = ?)`,
		libId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action that this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not get books from Book_detail. Error: %s", err)
		return err
	}
	defer resp.Close()

	for resp.Next() {
		title := ""
		bookKey := ""
		img := ""
		authorName := ""
		authorKey := ""
		err = resp.Scan(&title, &bookKey, &img, &authorName, &authorKey)
		if err != nil {
			ErrorRespone(c, `
                We could not perform this action at this moment. Please try again later.
                `, http.StatusInternalServerError)
			log.Printf("Could not scan for info in the query. Error: %s", err)
			return err
		}

		*myBooksPage = append(*myBooksPage, fmt.Sprintf(`
            <tr>
                <td class="book-display">
                    <div class="book-name">
                        <div class="book-img">
                            <a href="#" hx-post="/book"
                                hx-target=".contents"
                                hx-swap="innerHTML"
                                hx-vals='{
                                    "work"      :   "%s",
                                    "author"    :   "%s",
                                    "author_key":   "%s",
                                    "cover"     :   "%s"
                                }'
                                hx-push-url="true"
                            ><img src="%s" width="125px" height="200px">
                            </a>
                        </div>
                        <div class="book-title">
                            <h3> <a href="#"
                            hx-post="/book"
                            hx-target=".contentContainer"
                            hx-swap="innerHTML"
                            hx-vals='{
                                "work"      : "%s",
                                "author"    : "%s",
                                "author_key": "%s",
                                "cover"     : "%s"
                            }'
                            hx-trigger="click"
                            >%s</a></h3>
                            <p><a href="#"
                            hx-post="/author"
                            hx-target=".contentContainer"
                            hx-swap="innerHTML"
                            hx-vals='{
                                "key"       : "%s",
                                "bookKey"   : "%s",
                                "authorName": "%s"
                            }'>
                            %s</a></p>
                        </div>
                    </div>
                </td>
                <td class="actions">
                    <div class="btn-group" role="group"
                        style="max-height: 50px; max-width: 90%%; margin-left: -15px;">
                        <button type="button" class="btn btn-success firstOption"
                            style="width: 125px;">
                            <a hx-get="/move/finish"
                            hx-target=".myBookList"
                            hx-swap="innerHTML"
                            hx-vals='{
                                "key"   : "%s",
                                "from"  : "%s"
                                }'
                            style="font-size: 13px;"
                            >Finish Book</a>
                        </button>
                        <div class="dropdown bookBtn btn-group"
                            style="width: 5px;">
                            <button class="btn btn-success dropdown-toggle"
                                    type="button" id="wantToRead" data-bs-toggle="dropdown"
                                    aria-expanded="false"
                            >
                            </button>
                            <ul class="dropdown-menu">
                                <li><a class="dropdown-item"
                                    href="#"
                                    hx-get="/move/favorite"
                                    hx-target=".contents"
                                    hx-swap="innerHTML"
                                    hx-trigger="click"
                                    hx-vals='{
                                        "key"   : "%s",
                                        "from"  : "%s"
                                        }'
                                    >Move to Favorite</a></li>
                                <li><a class="dropdown-item" 
                                    href="#"
                                    hx-get="/move/toread"
                                    hx-target=".myBookList"
                                    hx-swap="innerHTML"
                                    hx-trigger="click"
                                    hx-vals='{
                                        "key"   : "%s",
                                        "from"  : "%s"
                                        }'
                                    hx-on::after-request="
                                        if (event.detail.xhr.status >= 400){
                                        document.querySelector('.responeMessage').innerHTML = event.detail.xhr.responseText;
                                        }"
                                    >Move to Plan to Read</a></li>
                                <li><a class="dropdown-item" 
                                    href="#"
                                    hx-get="/move/drop"
                                    hx-target=".contents"
                                    hx-swap="innerHTML"
                                    hx-trigger="click"
                                    hx-vals='{
                                        "key"   : "%s",
                                        "from"  : "%s"
                                        }'
                                    hx-on::after-request="
                                        if (event.detail.xhr.status >= 400){
                                        document.querySelector('.responeMessage').innerHTML = event.detail.xhr.responseText;
                                        }"
                                    >Drop Book</a></li>
                            </ul>
                    </div>
                </td>
            </tr>
            `, bookKey, authorName, authorKey, img, img,
			bookKey, authorName, authorKey, img, title,
			authorKey, bookKey, authorName, authorName,
			bookKey, "reading", bookKey, "reading",
			bookKey, "reading", bookKey, "reading"))
	}

	*myBooksPage = append(*myBooksPage, `
        </tbody>
        </table>
        </div>
        </div>
    `)
	return nil
}

func loadOptions(myPage *[]string) {
	*myPage = append(*myPage, fmt.Sprint(`
        <div class="contentContainer">
        <div class="contentLeft border rounded p-3 bg-light" style="max-width: 250px;" >
            <ul class="nav nav-pills flex-column">
            <li class="nav-item">
                <a class="nav-link reading active" 
                href="#"
                hx-get="/my-books/reading"
                hx-push-url="true"
                hx-trigger="click"
                hx-target=".myBookList"
                hx-trigger="click"
                hx-on::after-request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.reading').classList.add('active');
                    }
                "
                >Reading</a></li>
            <li class="nav-item">
                <a class="nav-link already" 
                href="#"
                hx-get="/my-books/alreadyRead"
                hx-target=".myBookList"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
                hx-on::after-Request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.already').classList.add('active');
                    }
                "
                >Already Read</a></li>
            <li class="nav-item">
                <a class="nav-link planning" 
                href="#"
                hx-get="/my-books/wantToRead"
                hx-target=".myBookList"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
                hx-on::after-Request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.planning').classList.add('active');
                    }
                "
            >Planning to Read</a></li>
            <li class="nav-item">
                <a class="nav-link favorite" 
                href="#" 
                hx-get="/my-books/favorites"
                hx-target=".myBookList"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
                hx-on::after-Request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.favorite').classList.add('active');
                    }
                "
            >Favorites</a></li>
            </ul>
            </div>
        <div class="contentRight" style="margin-left: auto;">
                <table class="table table-hover">
                    <thead>
                        <tr>
                            <th>Book Title</th>
                            <th>Edit</th>
                        </tr>
                    </thead>
                    <tbody class="myBookList">
    `))
}
