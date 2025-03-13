package search

import "github.com/gin-gonic/gin"

func SearchPage(c *gin.Context) {
	c.String(200, `
    <form hx-post="/book-search" hx-swap="innerHTML"
                hx-target=".search-display">
    <div class="search-box">
        <p>Search</p>
        <input type="text" autocomplete="off" name="query" placeholder="Book Title" required>
        <button type="submit">Search</button>
    </div>
    </form>
    <div class="search-display"></div>
`)
}
