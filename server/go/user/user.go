package user

import (
	"book/recomend"
	"fmt"
	"net/http"
	"strconv"

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

	if user.Fname == "" {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
		<div class="userProfile">
			<fieldset>
				<label for="first-name">First Name:
					<input id="first-name" type="text" placeholder="Enter your first name"/>
				</label>
				<label for="last-name">Last Name:
					<input id="last-name" type="text" placeholder="Enter your last name"/>
				</label>
				<label for="email">Email:
					<input id="email" type="text" placeholder="Enter your email"/>
				</label>
			</fieldset>
		</div>
	`)
	} else {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, fmt.Sprintf(`
			<div class="userProfile">
				<fieldset>
					<label for="first-name">First Name:
						<input id="first-name" type="text" value="%s"/>
					</label>
					<label for="last-name">Last Name:
						<input id="last-name" type="text" value="%s"/>
					</label>
					<label for="email">Email:
						<input id="email" type="text" value="%s"/>
					</label>
				</fieldset>
			</div>
		`, user.Fname, user.Lname, user.Email))
	}
}
