package user

import (
	"github.com/gin-gonic/gin"
)

func ChangePassBtn(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, `
	<form class="user-signin" 
			hx-post="/change-pass" 
			hx-on::after-request=" 
				if (event.detail.xhr.status >= 400) { 
					document.getElementById('errorWarning').innerHTML = event.detail.xhr.responseText; 
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
                <h3>Password Change</h3>
                    <label for="currPass">Current Password:</label><br>
                    <input type="password" id="currPass" name="currPass" 
                            required"
                            placeholder="Enter your current password">
                     <img src="https://media.geeksforgeeks.org/wp-content/uploads/20210917145551/eye.png" 
                        width="5%" height="5%" style="display:inlinel margin-left: -1.5%; vertical-align: middle"
                        id="currPass-enter"><br>
                    <label for="signup-password">New Password:</label><br>
                    <input type="password" id="signup-password" name="newPass" required
                            placeholder="Enter your password">
                     <img src="https://media.geeksforgeeks.org/wp-content/uploads/20210917145551/eye.png" 
                        width="5%" height="5%" style="display:inlinel margin-left: -1.5%; vertical-align: middle"
                        id="signup-enter">
                    <br> 
                    <p id="password-regex" style="color: red; display: none; font-size: 14px;">
                        Password need to be at least 1 character long with at least 1 number and a special character(!@#$%^/_)
                    </p>
                    <label for="password-reenter">Re-enter your password:</label><br>
                    <input type="password" id="password-reenter" required
                            placeholder="Reenter your password">
                     <img src="https://media.geeksforgeeks.org/wp-content/uploads/20210917145551/eye.png" 
                        width="5%" height="5%" style="display:inlinel margin-left: -1.5%; vertical-align: middle"
                        id="signup-reenter">
                    <br>
					<p id="reenter-regex" style="color: red; display: none; font-size: 14px;">
									Password need to be at least 10 character long with at least 1 number and a special character(!@#$%^/_)
                    </p>
        <p id="mismatch-pass" style="color:red; display: none; font-size: 14px;">Password does not match!</p><br>
        <div id="warning-msg">
        </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
            <button type="submit" id="submit-signup" class="btn btn-primary">Register</button>
          </div>
        </form>
		</div>
	  </div>
	</div>
	</form>
            `)

}
