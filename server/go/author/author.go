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
	Title string `json:"title"`
	Url   string `json:"url"`
}

type AuthorBio struct {
	Type  string
	Value string
}

func (ab *AuthorBio) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		*ab = AuthorBio{Value: str}
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
	*ab = AuthorBio{
		Type:  obj.Type,
		Value: obj.Value,
	}
	return nil
}

type Author struct {
	Key     string    `json:"key"`
	Bio     AuthorBio `json:"bio"`
	Birth   string    `json:"birth_date"`
	Death   string    `json:"death_date"`
	Links   []Link    `json:"links"`
	Name    string
	Photo   string
	BookKey string
}

func findAuthor(authorKey, authorName, bookKey string, db *sqlx.DB) (Author, error) {

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
		author.Name = authorName

		err = addAuthorToDB(author, db)
		if err != nil {
			return Author{}, err
		}

		return author, nil
	}
}

func lookUpAuthor(authorKey string) (Author, error) {
	url := "https://openlibrary.org"
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

	_, err := db.Exec(`INSERT INTO Author VALUES (?, ?, ?, ?, ?, ?)`,
		author.Key, author.Name, author.Birth, author.Photo, author.Death, author.Bio.Value)

	if err != nil {
		return err
	}

	for _, link := range author.Links {

		_, err = db.Exec(`INSERT INTO Links (title, url) VALUES (?, ?)`, link.Title, link.Url)
		if err != nil {
			return err
		}

		_, err = db.Exec(`INSERT INTO Author_links (author_id, link_id)
            SELECT ?, link_id FROM Links WHERE url = ?`, author.Key, link.Url)

		if err != nil {
			return err
		}
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

		author.Bio.Value = bio
		author.Key = key
		author.Name = name
		author.Birth = dob
		author.Death = dod
		author.Photo = photo
		err = getAuthorLinks(author, db)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to retrieved links for author. Error: %s", err))
		}
		return nil
	} else {
		return errors.New(fmt.Sprintf("Could not access the Author information with the key %s", authorKey))
	}
}

func getAuthorLinks(author *Author, db *sqlx.DB) error {
	ans, err := db.Query(`SELECT title, url FROM Links WHERE link_id IN 
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
