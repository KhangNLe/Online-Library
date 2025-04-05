package recomend

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetRecoommend(c *gin.Context, db *sqlx.DB, userId string) {
	user, err := strconv.Atoi(userId)
	if err != nil {
		log.Printf("Could not get user_id from %s. Error: %s",
			userId, err)
		ErrorRespone(c, `
            We could not respone to this request with your user id.
            Please contact the dev at the bottom of the page for fixes.
            `, http.StatusInternalServerError)
		return
	}
	var genres []string
	var book_ids []int
	err = getGenres(db, user, &genres)
	if err != nil {
		ErrorRespone(c, `
            We could not be able to collect your genres.
            Please contact the dev at the bottom of the page.
            `, http.StatusInternalServerError)
		return
	}
	err = getUserBookId(db, user, &book_ids)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action.
            Please contact the dev at the bottom of the page for fixes.
            `, http.StatusInternalServerError)
		return
	}

	topGenre := make(map[string]int)
	for _, genre := range genres {
		if _, ok := topGenre[genre]; ok {
			topGenre[genre]++
		} else {
			topGenre[genre] = 1
		}
	}
	var genre string
	amount := 0
	for key, val := range topGenre {
		if val > amount && key != "Fiction" {
			genre = key
			amount = val
		}
	}
	log.Println(genre)

	var books []string
	err = getBook(genre, &books, &book_ids)
	if err != nil {
		ErrorRespone(c, `
            We could not display your reccommend books.
            Please contact the dev at the bottom of the page.
            `, http.StatusInternalServerError)
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, strings.Join(books, "\n"))
}

func ErrorRespone(c *gin.Context, msg string, status int) {
	c.Header("Content-Type", "text/html")
	c.String(status, fmt.Sprintf(`
        <div class="errormsg" style="margin: auto;">
            <p style="font-size: 25px; color: red;>%s</p>
        </div>
    `, msg))
}

func getGenres(db *sqlx.DB, userId int, genres *[]string) error {
	var libId int
	resp, err := db.Query(`SELECT library_id FROM User_library 
            WHERE user_id = ?`, userId)
	if err != nil {
		log.Printf("Could not get the library_id. Error: %s", err)
		return err
	}
	defer resp.Close()

	if !resp.Next() {
		log.Printf("Could not find a library_id with user_id: %d", userId)
		return err
	}

	err = resp.Scan(&libId)
	if err != nil {
		log.Printf("Could not scan for library_id. Error: %s", err)
		return err
	}

	resp, err = db.Query(`SELECT G.genre_name
        FROM Genre G JOIN Book_Genre Bg
        ON G.genre_id = Bg.genre_id 
        WHERE Bg.book_id IN (
            SELECT book_id FROM Reading
            WHERE library_id = ?)
        OR Bg.book_id IN (
            SELECT book_id FROM Favorites 
            WHERE library_id  = ?)
        OR Bg.book_id IN (
            SELECT book_id FROM Read_Book 
            WHERE library_id = ?)
        `, libId, libId, libId)
	if err != nil {
		log.Printf("Could not get the genre. Error: %s", err)
		return err
	}
	defer resp.Close()

	for resp.Next() {
		var genre string
		err = resp.Scan(&genre)
		if err != nil {
			log.Printf("Could not scan for genre. Error: %s", err)
			return err
		}
		*genres = append(*genres, genre)
	}
	return nil
}

func getUserBookId(db *sqlx.DB, userId int, book_id *[]int) error {
	resp, err := db.Query(`SELECT book_id FROM Book_Detail 
        WHERE book_id IN (
            SELECT R.book_id FROM Reading R
            JOIN User_library U 
            ON R.library_id = U.library_id 
            WHERE U.user_id = ?
            )
        OR book_id IN (
            SELECT P.book_id FROM Planning_to_Read P
            JOIN User_library U
            ON P.library_id = U.library_id
            WHERE U.user_id = ?
            )
        OR book_id IN (
            SELECT F.book_id FROM Favorites F
            JOIN User_library U
            ON F.library_id = U.library_id
            WHERE U.user_id = ?
            )
        OR book_id IN (
            SELECT Rb.book_id FROM Read_Book Rb 
            JOIN User_library U
            ON Rb.library_id = U.library_id 
            WHERE U.user_id = ?
            )`, userId, userId, userId, userId)
	if err != nil {
		log.Printf("Could not collect book_ids from user library. Error: %s",
			err)
		return err
	}
	defer resp.Close()

	for resp.Next() {
		var key string
		err = resp.Scan(&key)
		if err != nil {
			log.Printf("Could not scan for book_id. Error: %s", err)
			return err
		}
		key = strings.Replace(key, "/works/OL", "", 1)
		key = strings.Replace(key, "W", "", 1)
		keyNum, err := strconv.Atoi(key)
		if err != nil {
			log.Printf("Could not convert %s to number", key)
			return err
		}
		*book_id = append(*book_id, keyNum)
	}
	return nil
}
