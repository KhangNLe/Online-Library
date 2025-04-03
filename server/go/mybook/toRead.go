package mybook

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func getToReadBooks(c *gin.Context, query *sqlx.DB, libId int,
	myPage *[]string) error {

	resp, err := query.Query(`SELECT Bd.book_id, Bd.cover_img, 
        Bd.title, A.author_id, A.name
        FROM Book B JOIN Author A ON B.author_id = A.author_id 
        JOIN Book_Detail Bd ON B.book_id = Bd.book_id 
        WHERE B.book_id IN (
            SELECT book_id FROM Planning_to_Read 
            WHERE library_id = ?)`, libId)
	if err != nil {
		ErrorRespone(c, ``, http.StatusBadRequest)
		log.Printf("Could not get info from planning_to_read. Error: %s", err)
		return err
	}
	defer resp.Close()

	for resp.Next() {
		book_id := ""
		img := ""
		title := ""
		name := ""
		author_id := ""
		err = resp.Scan(&book_id, &img, &title, &author_id, &name)
		if err != nil {
			ErrorRespone(c, ``, http.StatusInternalServerError)
			log.Printf("Could not scan for items in query. Error: %s", err)
			return err
		}

		*myPage = append(*myPage, fmt.Sprintf(`
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
                            <a hx-post="/my-books/reading"
                            hx-target=".myBookList"
                            hx-swap="innerHTML"
                            hx-vals='{
                                "key"   : "%s",
                                "from"  : "%s"
                                }'
                            style="font-size: 13px;"
                            >Reading</a>
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
                                    hx-post="/my-books/favorite"
                                    hx-target=".myBookList"
                                    hx-swap="innerHTML"
                                    hx-trigger="click"
                                    hx-vals='{
                                        "key"   : "%s",
                                        "from"  : "%s"
                                        }'
                                        hx-on::after-request="
                                            if (event.detail.xhr.status >= 400){
                                            document.querySelector('.myBookList').innerHTML = event.detail.xhr.responseText;
                                        }"
                                    >Move to Favorite</a></li>
                                <li><a class="dropdown-item" 
                                    href="#"
                                    hx-post="/my-books/drop"
                                    hx-target=".myBookList"
                                    hx-swap="innerHTML"
                                    hx-trigger="click"
                                    hx-vals='{
                                        "key"   : "%s",
                                        "from"  : "%s"
                                        }'
                                    hx-on::after-request="
                                        if (event.detail.xhr.status >= 400){
                                        document.querySelector('.myBookList').innerHTML = event.detail.xhr.responseText;
                                        }"
                                    >Drop Book</a></li>
                            </ul>
                    </div>
                </td>
            </tr>
        `, book_id, name, author_id, img, img,
			book_id, name, author_id, img, title,
			author_id, book_id, name, name,
			book_id, "toread", book_id, "toread",
			book_id, "toread"))
	}

	*myPage = append(*myPage, `
        </tbody>
        </table>
        </div>
        </div>
    `)
	return nil
}
