package user

import (
	"book/mybook"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EditProfileBtn(c *gin.Context) {

	userDetail := make(map[string]string)
	err := c.Bind(&userDetail)
	if err != nil {
		mybook.ErrorRespone(c, `
			Could not perform this action at the moment. Please contact the dev for this problem.
			`, http.StatusInsufficientStorage)
		return
	}

	fname, ok := userDetail["fname"]
	if !ok {
		fname = ""
		log.Printf("Could not find fname in c.Bind")
	}
	lname, ok := userDetail["lname"]
	if !ok {
		fname = ""
		log.Println("Could not find lname in c.Bind")
	}
	email, ok := userDetail["email"]
	if !ok {
		email = ""
		log.Println("Could not find email in c.Bind")
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf(`
		<form class="user-signin" 
				hx-post="/edit-profile" 
				hx-target=".contents"
				hx-swap="innerHTML"
				hx-on::after-request=" 
					if (event.detail.xhr.status >= 400) { 
						document.getElementById('warning-msg').innerHTML = event.detail.xhr.responseText; 
					} else if (event.detail.xhr.status === 200){
						document.querySelector('.btn-close').click();	
					}
				"> 
		<div class="modal fade" id="staticBackdrop"
				data-bs-backdrop="static" 
				data-bs-keyboard="true" 
				tabindex="-1" 
				aria-labelledby="staticBackdropLabel" 
				aria-hidden="false">
		  <div class="modal-dialog">
			<div class="modal-content">
			  <div class="modal-header">
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			  </div>
				<div class="body-footer">
			  <div class="modal-body">
					<h3>Edit Profile</h3>
						<label for="currPass">First Name:</label><br>
						<input type="text" id="currPass" name="fname" 
							autocomplete="off"
							placeholder="Enter your first name" value="%s">
						 <br>
						<label for="signup-password">Last Name:</label><br>
						<input type="text" id="signup-password" name="lname"
							autocomplete="off"
							placeholder="Enter your last name" value="%s">
						<br> 
						<label for="password-reenter">Email:</label><br>
						<input type="text" id="password-reenter" name="email"
							autocomplete="off"
							placeholder="Enter your email" value="%s">
						<br>
			<p id="mismatch-pass" style="color:red; display: none; font-size: 14px;">Password does not match!</p><br>
			<div id="warning-msg">
			</div>
			  </div>
			  <div class="modal-footer">
				<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
				<button type="submit" id="profile-change" class="btn btn-primary">Submit Change</button>
			  </div>
			</form>
			</div>
		  </div>
		</div>
		</form>
		`, fname, lname, email))
}
