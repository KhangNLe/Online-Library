package search

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
)

type Works struct {
	Results []SubjectResult `json:"works"`
}

type Author struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type SubjectResult struct {
	Key     string   `json:"key"`
	Title   string   `json:"title"`
	IMG     int      `json:"cover_id"`
	Authors []Author `json:"authors"`
}

var subjectBooks []SubjectResult

func SearchSubject(text string) []SubjectResult {
	url := "https://openlibrary.org/subjects/"
	field := ".json?limit=250"

	return getSearchSubject(url + text + field)
}

func getSearchSubject(searchStr string) []SubjectResult {

	var search Works

	resp, err := http.Get(searchStr)
	if err != nil {
		log.Printf("Could not fetch the %s url, erro: %s", (searchStr), err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		log.Printf("Could not decode the api, err: %s", err)
	}

	subjectBooks = search.Results

	return search.Results
}

func appendSubjectDisplay(start, end, totalBook int, bookDisplay *[]string) {

	for i := start; i < end; i++ {
		if i == totalBook {
			break
		}

		if len(subjectBooks[i].Authors) == 0 {
			continue
		}
		var bookPic string
		if subjectBooks[i].IMG == 0 {
			bookPic = "https://upload.wikimedia.org/wikipedia/commons/1/14/No_Image_Available.jpg?20200913095930"
		} else {
			bookPic = html.EscapeString("https://covers.openlibrary.org/b/id/" + strconv.Itoa(subjectBooks[i].IMG) + "-M.jpg")
		}

		*bookDisplay = append(*bookDisplay,
			fmt.Sprintf(`
                <div class="books">
                <img src="%s">
                    <a hx-post="/book" hx-swap="innerHTML"
                        hx-trigger="click"
                        hx-target=".contents" 
                        hx-vals='{
                            "work":     "%s",
                            "author":   "%s",
                            "author_key":   "%s",
                            "cover":    "%s"}'
                        hx-replace-url="/book%s"
                        href="#"
                    >
                    %s</a>
                </div>`,
				bookPic,
				html.EscapeString(subjectBooks[i].Key),
				subjectBooks[i].Authors[0].Name,
				subjectBooks[i].Authors[0].Key,
				bookPic,
				subjectBooks[i].Key,
				subjectBooks[i].Title,
			),
		)
	}
}
