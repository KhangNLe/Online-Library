package search

import "github.com/gin-gonic/gin"

func SearchPage(c *gin.Context) {

	c.Header("Content-Type", "text/html")
	c.String(200, `
        <nav class="navbar bg-body-tertiary">
            <div class="container-fluid" style="max-width: fit-content; margin-left: auto; margin-right:auto;">
                <form class="d-flex" role="search"
                        hx-post="/book-search"
                        hx-swap="innerHTML"
                        hx-target=".display"
                        hx-vals='{"page" : "1"}'
                        hx-push-url="true"
                        hx-on::after-request=" if (event.detail.xhr.status >= 400) { document.querySelector('.search-display').innerHTML = event.detail.xhr.responseText; }" 
                        >
            <input class="form-control me-2" size="75%" name="query" type="search" autocomplete="off" placeholder="Enter title or author" aria-label="Search">
                    <button id="searchBtn" class="btn btn-outline-success" type="submit">Search</button>
                </form>
            </div>
        </nav>
        <div class="display"></div>
        `)
}
