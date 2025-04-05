package move

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AddingToLibrary(userId string, c *gin.Context, db *sqlx.DB, optinon int) {
	hxVals := make(map[string]string)

	err := c.Bind(&hxVals)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("There are no hx-vals on button. Error: %s", err)
		return
	}

	key, _ := hxVals["key"]

	query, err := db.Beginx()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Error starting db.Beginx. Error: %s", err)
		return
	}
	defer func() {
		if err != nil {
			query.Rollback()
		}
	}()

	num, err := strconv.Atoi(userId)
	if err != nil {
		ErrorRespone(c, ``, http.StatusInternalServerError)
		log.Printf("Couldn't convert userID (%s) back to int. Error: %s",
			userId, err)
	}

	libraryId, err := query.Query(`SELECT library_id FROM User_library WHERE user_id = ?`, num)
	if err != nil {
		ErrorRespone(c, `
            Current could not access the database at the moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Error with library_id search. Error: %s", err)
		return
	}
	defer libraryId.Close()

	var libID int

	if !libraryId.Next() {
		ErrorRespone(c, `
            We could not find your library. Please contact support.
            `, http.StatusInternalServerError)
		log.Println("There is no library_id for this user")
		return
	}
	err = libraryId.Scan(&libID)
	if err != nil {
		ErrorRespone(c, `
            We could not complete this action at the moment, please try again later.
            `, http.StatusBadRequest)
		log.Printf("Something wrong why trying to scan for library_id. Error: %s", err)
		return
	}

	libSesh := ""
	bookOrAuthor := ""
	switch optinon {
	case 0:
		err = wantToRead(c, query, libID, key, num)
		libSesh = "Wanting to Read"
		bookOrAuthor = "book"
	case 1:
		err = addToReading(c, query, libID, key, num)
		libSesh = "Reading"
		bookOrAuthor = "book"
	case 2:
		err = addFavoriteAuthor(c, query, key, libID)
		libSesh = "Favorite Author"
		bookOrAuthor = "author"
	case 3:
		err = addToBlockAuthor(c, query, libID, key)
		libSesh = "Block Author"
		bookOrAuthor = "author"
	case 4:
		err = addToFavoriteBook(c, query, libID, key)
		libSesh = "Favorite Book"
		bookOrAuthor = "book"
	case 69:
		err = addToAlreadyRead(c, query, libID, key, num)
		libSesh = "Already Read"
		bookOrAuthor = "book"
	default:
		err = errors.New(fmt.Sprintf(`
            The %d not match any of the option.
        `, optinon))
	}

	if err != nil {
		log.Printf("Error happened when tried to add book. Error: %s", err)
		return
	}

	err = query.Commit()
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment, please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not commit change to db. Error: %s", err)
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf(`
        <br><p style="color: green; font-size: 15px; margin-left: -20px;">
        The %s is added to your %s session.
        </p></br>
    `, bookOrAuthor, libSesh))

}

func ErrorRespone(c *gin.Context, msg string, status int) {
	c.Header("Content-Type", "text/html")
	c.String(status, fmt.Sprintf(`
        <br><p style="color: red; font-size: 15px; margin-left: -20px;">
        %s
        </p></br>
        `, msg))
}
