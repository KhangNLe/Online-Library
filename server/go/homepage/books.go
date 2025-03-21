package homepage

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ResultGuten struct {
	Results []BookGuten `json:"results"`
}

type BookGuten struct {
	Title string `json:"title"`
}

type ResultBook struct {
	Results string `json:"docs"`
}

type Book struct {
	Key        string   `json:"key"`
	Title      string   `json:"title"`
	IMG        int      `json:"cover_i"`
	AuthorKey  []string `json:"author_key"`
	AuthorName []string `json:"author_name"`
}

func TopFiction() ([]Book, error) {
	urlGuten := "https://gutendex.com/books/?topic=fiction&sort=popular"
	var popular ResultGuten

	resp, err := http.Get(urlGuten)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&popular)
	if err != nil {
		return nil, err
	}
	var popularBook []Book
	for idx, result := range popular.Results {
		if idx == 25 {
			break
		}
		book, err := getBook(result)
		if err != nil {
			idx--
			continue
		}
		popularBook = append(popularBook, book)

	}

	return popularBook, nil
}

func getBook(info BookGuten) (Book, error) {
	url := "https://openlibrary.org/search.json?q="
	features := "&fields=author_key,cover_i,title,key"
	info.Title = strings.ReplaceAll(info.Title, " ", "+")

	var book Book
	resp, err := http.Get(url + info.Title + features)
	if err != nil {
		return Book{}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&book)
	if err != nil {
		return Book{}, err
	}

	return book, nil
}
