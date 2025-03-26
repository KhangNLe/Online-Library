package author

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Link struct {
	Title string `json:"tittle"`
	Url   string `json:"url"`
}

type Author struct {
	Key     string `json:"key"`
	Bio     string `json:"bio"`
	Name    string `json:"personal_name"`
	Birth   string `json:"birth_date"`
	Death   string `json:"death_date"`
	Links   []Link `json:"links"`
	Photo   string
	BookKey string
}

func findAuthor(authorKey, bookKey string, db *sqlx.DB) (Author, error) {

	ok, err := isExist(authorKey, db)
	if err != nil {
		return Author{}, err
	}
	if ok {
		var author Author
		err = getAuthorInfo(&author, authorKey, db)
		if err != nil {
			return Author{}, err
		}
		return author, nil
	} else {
		author, err := lookUpAuthor(authorKey)
		if err != nil {
			return Author{}, err
		}

		key := strings.Replace(authorKey, "/authors/", "", 1)
		photoUrl := "https://covers.openlibrary.org/a/olid/" + key + "-M.jpg"
		author.Photo = photoUrl
		author.BookKey = bookKey

		err = addAuthorToDB(author, db)
		if err != nil {

		}

		return author, nil
	}
}

func lookUpAuthor(authorKey string) (Author, error) {
	url := "https://openlibrary.org/"
	tail := ".json"

	resp, err := http.Get(url + authorKey + tail)
	if err != nil {
		return Author{}, err
	}
	defer resp.Body.Close()

	var author Author
	err = json.NewDecoder(resp.Body).Decode(&author)
	if err != nil {
		return Author{}, err
	}

	return author, nil
}

func addAuthorToDB(author Author, db *sqlx.DB) error {

	if author.Death == "" {
		author.Death = "Still alive"
	}
	if author.Birth == "" {
		author.Birth = "Unknown birth date"
	}
	_, err := db.Exec(`INSERT INTO Author VALUES (?, ?, ?, ?, ?, ?, ?)`,
		author.Key, author.Name, author.Birth, author.Photo, author.Death, author.Bio)

	if err != nil {
		return err
	}

	for _, link := range author.Links {

		_, err = db.Exec(`INSERT INTO Links VALUES (?, ?)`, link.Title, link.Url)
		if err != nil {
			return err
		}

		_, err = db.Exec(`INSERT INTO Author_links VALUES (author_id, link_id)
            SELECT ?, link_id FROM Links WHERE url = ?)`, author.Key, link.Url)

		if err != nil {
			return err
		}
	}

	_, err = db.Exec(`INSERT INTO Book VALUES (?, ?)`, author.BookKey, author.Key)
	if err != nil {
		return err
	}

	return nil
}

func isExist(authorKey string, db *sqlx.DB) (bool, error) {
	ans, err := db.Query(`SELECT * FROM Author WHERE author_id = ?`, authorKey)
	if err != nil {
		return false, err
	}

	defer ans.Close()

	return ans.Next(), nil
}

func getAuthorInfo(author *Author, authorKey string, db *sqlx.DB) error {
	key := ""
	name := ""
	bio := ""
	dob := ""
	dod := ""
	photo := ""

	ans, err := db.Query(`SELECT * FROM Author WHERE author_id = ?`, authorKey)
	if err != nil {
		return err
	}
	defer ans.Close()

	if ans.Next() {
		err = ans.Scan(&key, &name, &dob, &photo, &dod, &bio)
		if err != nil {
			return err
		}

		author.Bio = bio
		author.Key = key
		author.Name = name
		author.Birth = dob
		author.Death = dod
		author.Photo = photo
		return nil
	} else {
		return errors.New(fmt.Sprintf("Could not access the Author information with the key %s", authorKey))
	}
}

func getAuthorLinks(author *Author, db *sqlx.DB) error {
	ans, err := db.Query(`SELECT title, url FROM Links WHERE link_id = 
            (SELECT link_id FROM Author_links WHERE author_id = ?)`, author.Key)

	if err != nil {
		return err
	}
	defer ans.Close()

	for ans.Next() {
		title := ""
		url := ""
		err = ans.Scan(&title, &url)
		if err != nil {
			return err
		}
		author.Links = append(author.Links, Link{Title: title, Url: url})
	}
	return nil
}
