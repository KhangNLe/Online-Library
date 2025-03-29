package author

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetAuthor(c *gin.Context, db *sqlx.DB) {
	booksVals := make(map[string]string)

	err := c.Bind(&booksVals)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Could not get the info from bookDetail. Error: %s", err)
		return
	}
	authorKey, ok := booksVals["key"]
	if !ok {
		c.Status(http.StatusInternalServerError)
		log.Println("Could not find an author key")
		return
	}

	bookKey, ok := booksVals["bookKey"]
	if !ok {
		c.Status(http.StatusInternalServerError)
		log.Println("Unable to get the book key for this author")
		return
	}

	author, err := findAuthor(authorKey, bookKey, db)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Error while looking up author. Error: %s", err)
		return
	}

	printAuthor(c, author)
}

func printAuthor(c *gin.Context, author Author) {
	c.Header("Content-Type", "text/html")
	var authorPage []string
	authorPage = append(authorPage, fmt.Sprintf(`
            <div class="contentLeft">
                <div class="bookImg">
                    <img src="%s">
                </div>
                <div class="bookAction">
                    <div class="btn-group" role="group" 
                        style="max-height: 60px; margin-left: -12%% ;">
                        <button type="button" class="btn btn-success"
                            style="width: 250px;">
                            <a hx-get="/favorite-author/add"
                            hx-target=".responeMessage"
                            hx-swap="innerHTML"
                            hx-vals='{
                                "authorKey": "%s"
                                }'
                            hx-on::after-request="
                                if (event.detail.xhr.status >= 400){
                                    document.querySelector('.responeMessage').innerHTML = event.detail.xhr.responseText;
                                }"
                            >Add to Favorite</a>
                        </button>
                            <div class="dropdown btn-group" style="width: 30px;">
                            <button class="btn btn-success dropdown-toggle"
                                    type="button" data-bs-toggle="dropdown"
                                    aria-expanded="false"
                            ></button>
                            <ul class="dropdown-menu">
                                <li><a class="dropdown-item" href="#">Block Author</a></li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
            <div class="contentRight">
                <div class="bookTitle">
        <p style="font-size: 25px;">Author: %s</p>
                </div>
                <div class="dob">
                    <span>
                        <h3 style="font-size: 17px; display: inline;">Birth Date:</h3>
                        <p style="display: inline;">    %s</p>
                    </spane>
                </div>
                <div class="dod">
                    <span>
                        <h3 style="font-size: 17px; display: inline;">Death Date:</h3>
                        <p style="display: inline;">    %s</p>
                    </span>
                </div>  
                <div class="bookDescription" style="max-width: 75%%;">
                    <h3>Bio:</h3>
                    <p>%s</p>
                </div>
    `, author.Photo,
		author.Key,
		author.Name,
		author.Birth,
		author.Death,
		author.Bio))
	linksDisplay(&authorPage, author)
	c.String(200, strings.Join(authorPage, ""))
}

func linksDisplay(authorPage *[]string, author Author) {
	*authorPage = append(*authorPage, fmt.Sprintf(`
        <div class="accordion accordion-flush" id="accordionFlushLinks"
            style="max-width: 275px;">
          <div class="accordion-item">
            <h2 class="accordion-header">
              <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#flush-collapseOne" aria-expanded="false" aria-controls="flush-collapseOne">
                Author's Links
              </button>
            </h2>
            <div id="flush-collapseOne" class="accordion-collapse collapse" data-bs-parent="#accordionFlushLinks">
                <div class="accordion accordion-flush" id="accordionFlushTitle">
        `))
	for idx, link := range author.Links {
		str := strconv.Itoa(idx)
		*authorPage = append(*authorPage, fmt.Sprintf(`
            <div class="accordion-item">
                <h2 class="accordion-header-link">
                    <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="%s" aria-expanded="false" aria-controls="flush-collapseOne">
                        %s
                    </button>
                </h2>
            <div id="%s" class="accordion-collapse collapse" data-bs-parent="#accordionFlushTitle">
                <div class="accordion-body"><p><a href="%s">%s</a></p></div>
                </div>
            </div>
			`, ("#flush-collapse"+str),
			link.Title,
			("flush-collapse"+str),
			link.Url,
			link.Title))
	}
	*authorPage = append(*authorPage, fmt.Sprintf(`
        </div>
        </div>
        </div>
        </div>
        </div>
        `))
}
