package user

import (
	"book/recomend"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func UserProfile(c *gin.Context, db *sqlx.DB, userId string) {
	user_num, err := strconv.Atoi(userId)
	if err != nil {
		recomend.ErrorRespone(c, `
			We could not get your user profile at this moment.
			Please contact the dev below for support.
			`, http.StatusInternalServerError)
		return
	}

	user, err := getUserProfile(db, user_num)
	if err != nil {
		recomend.ErrorRespone(c, `
			We could not get your user profile at this moment.
			Please contact the dev at the bottom of the page for support.
			`, http.StatusInternalServerError)
		return
	}
	var profile []string
	loadUserTable(&profile)

	profile = append(profile, `
		<form id="profile-change"
			hx-post="/change-profile"
			hx-target=".contents"
			hx-swap="innerHTML"
		>
	`)
	if user.Fname == "" {
		profile = append(profile, `
			<td class="user-display">
				<fieldset>
					<label for="first-name">First Name:
						<input id="first-name" type="text" placeholder="Enter your first name"/>
					</label>
				</fieldset>
			</td>
			</tr>
			<tr>
			<td class="user-display">
				<fieldset>
					<label for="last-name">Last Name:
						<input id="last-name" type="text" placeholder="Enter your last name"/>
					</label>
				</fieldset>
			</td>
			</tr>
			<tr>
			<td class="user-display">
				<fieldset>
					<label for="email">Email:
						<input id="email" type="text" placeholder="Enter your email"/>
					</label>
				</fieldset>
			</td>
			</tr>
	`)
	} else {
		profile = append(profile, fmt.Sprintf(`
			<td class="user-display">
				<fieldset>
					<label for="first-name">First Name:
						<input id="first-name" type="text" value="%s"/>
					</label>
				</fieldset>
			</td>
			</tr>
			<tr>
			<td class="user-display">
				<fieldset>
					<label for="last-name">Last Name:
						<input id="last-name" type="text" value="%s"/>
					</label>
				</fieldset>
			</td>
			</tr>
			<tr>
			<td class="user-display">
				<fieldset>
					<label for="email">Email:
						<input id="email" type="text" value="%s"/>
					</label>
				</fieldset>
			</td
			</tr>
		`, user.Fname, user.Lname, user.Email))
	}

	profile = append(profile, `
		<tr><td>
		<button type="submit" class="btn btn-success">
		Submit Change
		</button>
		</td></tr>
		</form>
		</tbody>
		</table>
		</div>
		</div>
	`)
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, strings.Join(profile, "\n"))
}

func loadUserTable(profile *[]string) {
	*profile = append(*profile, fmt.Sprintf(`
        <div class="contentContainer">
        <div class="contentLeft border rounded p-3 bg-light" 
            style="max-width: 250px; max-height: 250px; margin-top: 19px;">
            <ul class="nav nav-pills flex-column">
            <li class="nav-item">
                <a class="nav-link reading" 
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
						document.querySelector('.table').classList.remove('profile');
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
						const table = document.querySelector('.table');
						table.classList.remove('profile');
						table.classList.add('table-hover');
						document.querySelector('.tableHead').innerHTML = '<th>Book Title</th><th>Edit</th>';
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
						const table = document.querySelector('.table');
						table.classList.remove('profile');
						table.classList.add('table-hover');
						document.querySelector('.tableHead').innerHTML = '<th>Book Title</th><th>Edit</th>';
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
						const table = document.querySelector('.table');
						table.classList.remove('profile');
						table.classList.add('table-hover');
						document.querySelector('.tableHead').innerHTML = '<th>Book Title</th><th>Edit</th>';
                    }
                "
            >Favorites</a></li>
            <li class="nav-item">
                <a class="nav-link profile active" 
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
						const table = document.querySelector('.table');
                    }
                "
            >Profile</a></li>
            </ul>
            </div>
        <div class="contentRight" style="margin-left: auto;">
                <table class="table profile">
                    <thead>
                        <tr class="tableHead">
                            <th>Profile</th>
                        </tr>
                    </thead>
                    <tbody class="myBookList">
    `))
}
