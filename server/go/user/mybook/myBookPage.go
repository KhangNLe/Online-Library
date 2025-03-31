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
        <div class="contentLeft border rounded p-3 bg-light" style="max-width: 250px;" >
            <ul class="nav nav-pills flex-column">
            <li class="nav-item">
                <a class="nav-link reading active" 
                href="#"
                hx-get="/my-books/reading"
                hx-push-url="true"
                hx-trigger="click"
                hx-target=".contents"
                hx-trigger="click"
                hx-vals='{
                    "userId": "%s"
                }'
                hx-on::after-request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.reading').classList.add('active');
                    }
                "
                >Reading</a></li>
            <li class="nav-item">
                <a class="nav-link already" 
                href="#"
                hx-get="/my-books/alreadyRead"
                hx-target=".contentRight"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
                hx-vals='{
                    "userId": "%s"
                    }'
                hx-on::after-Request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.already').classList.add('active');
                    }
                "
                >Already Read</a></li>
            <li class="nav-item">
                <a class="nav-link planning" 
                href="#"
                hx-get="/my-books/wantToRead"
                hx-target=".contentRight"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
                hx-vals='{
                    "userId": "%s"
                    }'
                hx-on::after-Request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.planning').classList.add('active');
                    }
                "
            >Planning to Read</a></li>
            <li class="nav-item">
                <a class="nav-link favorite" 
                href="#" 
                hx-get="/my-books/favorites"
                hx-target=".contentRight"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
                hx-vals='{
                    "userId": "%s"
                    }'
                hx-on::after-Request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.favorite').classList.add('active');
                    }
                "
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

func ErrorRespone(c *gin.Context, msg string, status int) {
	c.Header("Content-Type", "text/html")
	c.String(status, fmt.Sprintf(`
        <br><p style="color: red; font-size: 15px;">
        %s
        </p></br>
        `, msg))
}
