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

	if user.Fname == "" {
		profile = append(profile, `
			<td class="user-display">
				<form name="change_fname"
				hx-post="/my-books/profile"
				hx-target=".contents"
				hx-swap="innerHTML"
				hx-on::after-request=" if (event.detail.xhr.status >= 400) {
					document.querySelector('.errorMsg').innerHTML = event.detail.xhr.responseText;
				} else if (event.detail.xhr.status == 200){
					alert('Profile added successfully');
				}"
				>
				<fieldset>
					<label for="first-name">First Name:
						<input id="first-name" name="fname" type="text"
							placeholder="Enter your first name"
							autocomplete="off"/>
					</label>
				<fieldset>
					<label for="last-name">Last Name:
						<input id="last-name" name="lname" type="text"
							placeholder="Enter your last name"
							autocomplete="off"/>
					</label>
				</fieldset>
				<fieldset>
					<label for="email">Email:
						<input id="email" type="text" name="email" 
							placeholder="Enter your email"
							autocomplete="off"/>
					</label>
				</fieldset>
				<button type="submit" class="btn btn-success"
				>Change Profile</button>
				</form>
				<button type="button" class="btn btn-success"
					hx-get="/change-pass"
					hx-target="#passwordChange"
					hx-swap="innerHTML"
				>Change Password</button>
			</td>
			</tr>
	`)
	} else {
		profile = append(profile, fmt.Sprintf(`
			<td class="user-display">
				<form name="change_fname"
				hx-target=".contents"
				hx-swap="innerHTML"
				hx-post="/my-books/profile"
				hx-on::after-request=" if (event.detail.xhr.status >= 400) {
					document.querySelector('.errorMsg').innerHTML = event.detail.xhr.responseText;
				} else if (event.detail.xhr.status === 200){
					alert('Change Request successed');
					console.log('batman');
				}"
				>
				<fieldset>
					<label for="first-name">First Name:
						<input id="first-name" name="fname" type="text"
							value="%s" autocomplete="off"/>
					</label>
				</fieldset>
				<fieldset>
					<label for="last-name">Last Name:
						<input id="last-name" name="lname" type="text"
							value="%s" autocomplete="off"/>
					</label>
				</fieldset>
				<fieldset>
					<label for="email">Email:
						<input id="email" type="text" name="email" 
							value="%s" autocomplete="off"/>
					</label>
				</fieldset>
				<button type="submit" class="btn btn-success"
				>Change Profile</button>
				</form>
				<div class="profileBtns">
					<button type="button" class="btn btn-success"
						hx-get="/change-pass"
						hx-target="#passwordChange"
						hx-swap="innerHTML"
						hx-trigger="click"
					>Change Password</button>
				</div>
			</td>
			</tr>
		`, user.Fname, user.Lname, user.Email))
	}

	profile = append(profile, `
		</tbody>
		</table>
			<div class="errorMsg"></div>
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
