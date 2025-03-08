package htmxswap

import "github.com/gin-gonic/gin"

func AboutPage(c *gin.Context) {
	c.Header("Context-Type", "text/html")
	c.String(200, `
    <h2>About</h2>       
        <p></p>
        `)
}
