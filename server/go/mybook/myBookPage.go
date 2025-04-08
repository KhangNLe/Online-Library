package mybook

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func MyBookPage(c *gin.Context, db *sqlx.DB, user string, dst string) {
	userId, err := strconv.Atoi(user)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at this moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Could not convert user_id to int. Error: %s", err)
		return
	}

	resp, err := db.Query(`SELECT library_id FROM User_library WHERE
                user_id = ?`, userId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Error when getting lib_id. Error: %s", err)
		return
	}
	defer resp.Close()

	if !resp.Next() {
		ErrorRespone(c, `
            Could not find library for this current user. Please contact the dev for this fix.
            `, http.StatusBadRequest)
		log.Printf("No library for the user_id %d", userId)
		return
	}

	var libId int
	err = resp.Scan(&libId)
	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment. Please try again later.
            `, http.StatusBadRequest)
		log.Printf("Could not scan for lib_id. Error: %s", err)
		return
	}

	var myPage []string

	switch dst {
	case "reading":
		err = currentlyReading(c, db, libId, &myPage)
	case "favorites":
		err = getFavoriteBook(c, db, libId, &myPage)
	case "toread":
		err = getToReadBooks(c, db, libId, &myPage)
	case "finish":
		err = getMyReadBook(c, db, libId, &myPage)
	default:
		err = errors.New(
			fmt.Sprintf("Could not find any matching option with %s", dst))
	}

	if err != nil {
		ErrorRespone(c, `
            We could not perform this action at the moment. Please try again later.
            `, http.StatusInternalServerError)
		log.Printf("Something wrong while trying to display book. Error: %s", err)
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, strings.Join(myPage, "\n"))
}

func ErrorRespone(c *gin.Context, msg string, status int) {
	c.Header("Content-Type", "text/html")
	c.String(status, fmt.Sprintf(`
        <br><p style="color: red; font-size: 15px;">
        %s
        </p></br>
        `, msg))
}

func loadOptions(myPage *[]string) {
	*myPage = append(*myPage, fmt.Sprint(`
        <div class="contentContainer">
        <div class="contentLeft border rounded p-3 bg-light" 
            style="max-width: 250px; max-height: 250px; margin-top: 19px;">
            <ul class="nav nav-pills flex-column">
            <li class="nav-item">
                <a class="nav-link reading active" 
                href="#"
                hx-get="/my-books/reading"
                hx-push-url="true"
                hx-trigger="click"
                hx-target=".contents"
                hx-trigger="click"
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
                hx-get="/my-books/finish"
                hx-target=".myBookList"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
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
                hx-get="/my-books/toread"
                hx-target=".myBookList"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
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
                hx-target=".myBookList"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
                hx-on::after-Request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.favorite').classList.add('active');
                    }
                "
            >Favorites</a></li>
            <li class="nav-item">
                <a class="nav-link profile" 
                href="#" 
                hx-get="/my-books/profile"
                hx-target=".contents"
                hx-swap="innerHTML"
                hx-push-url="true"
                hx-trigger="click"
                hx-on::after-Request= "
                    if (event.detail.xhr.status == 200) {
                        document.querySelectorAll('.nav-link').forEach(elmt => {
                            elmt.classList.remove('active');
                        });
                        document.querySelector('.profile').classList.add('active');
                    }
                "
            >Profile</a></li>
            </ul>
            </div>
        <div class="contentRight" style="margin-left: auto;">
                <table class="table table-hover">
                    <thead>
                        <tr>
                            <th>Book Title</th>
                            <th>Edit</th>
                        </tr>
                    </thead>
                    <tbody class="myBookList">
    `))
}
