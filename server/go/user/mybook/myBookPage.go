package mybook

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MyBookPage(c *gin.Context, user string) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf(`
        <div class="contentContainer">
         <div class="contentLeft">
            <ul class="nav nav-pills nav-stacked">
            <li class="active"><a 
                href="#"
                HX-Redirect="/my-books"
                >Reading</a></li>
            <li><a hx-get="/alreadyRead"
                hx-target=".contentRight"
                hx-swap="innerHTML"
                hx-trigger="click"
                hx-vals='{
                    "userId": "%s"
                    }'
                >Already Read</a></li>
            <li><a hx-get="/reading"
                hx-target=".contentRight"
                hx-swap="innerHTML"
                hx-trigger="click"
                hx-vals='{
                    "userId": "%s"
                    }'
                >Reading</a></li>
            <li><a hx-get="/wantToRead"
                hx-target=".contentsRight"
                hx-swap="innerHTML"
                hx-trigger="click"
                hx-vals='{
                    "userId": "%s"
                    }'
            >Planning to Read</a></li>
            <li><a hx-get="/favorites"
                hx-target=".contentsRight"
                hx-swap="innerHTML"
                hx-trigger="click"
                hx-vals='{
                    "userId": "%s"
                    }'
            >Favorites</a></li>
            </ul>
            </div>
            <div class="contentRight">
            </div>
        </div>
    `, user,
		user,
		user,
		user))
}
