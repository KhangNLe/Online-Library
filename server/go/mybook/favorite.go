package mybook

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func getFavoriteBook(c *gin.Context, query *sqlx.DB,
	libId int, myPage *[]string) error {

	resp, err := query.Query(`SELECT Bd.cover_img, Bd.book_id, 
        Bd.title, A.name, A.author_id
        FROM Book B JOIN Author A ON B.author_id = A.author_id
        JOIN Book_Detail Bd ON B.book_id = Bd.Book_id 
        WHERE B.book_id IN (SELECT book_id FROM Favorite_Book
                            WHERE library_id = ?)`, libId)
	if err != nil {
		ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Could not get book from Favorite_book. Error: %s", err)
		return err
	}
	defer resp.Close()

	for resp.Next() {
		img := ""
		bookKey := ""
		title := ""
		name := ""
		author_id := ""
		err = resp.Scan(&img, &bookKey, &title, &name, &author_id)
		if err != nil {
			ErrorRespone(c, ``, http.StatusBadRequest)
			log.Printf("Could not scan for the item in favorite_book. Error: %s", err)
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
                            <a hx-get="/move/favorite"
                            hx-target=".contents"
                            hx-swap="innerHTML"
                            hx-vals='{
                                "key"   : "%s",
                                "from"  : "%s"
                                }'
                            style="font-size: 13px;"
                            >Drop Book</a>
                        </button>
                </td>
            </tr>
            `, bookKey, name, author_id, img, img,
			bookKey, author_id, author_id, img, title,
			author_id, bookKey, author_id, name,
			bookKey, "favorite"))
	}

	*myPage = append(*myPage, `
        </tbody>
        </table>
        </div>
        </div>
    `)
	return nil
}
