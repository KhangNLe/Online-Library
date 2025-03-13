package htmxswap

import "github.com/gin-gonic/gin"

func SignUpBtn(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, `
        <form id="register-form" hx-post="/register"  hx-swap="innerHTML" hx-target="#warning-msg"
                hx-on::after-request="if (event.detail.xhr.status >= 400) { document.getElementById('warning-msg').innerHTML = event.detail.xhr.responseText; }">
          <div class="modal-body">
                <h3>Register</h3>
                    <label for="signup-userid">Username:</label><br>
                    <input type="text" id="signup-userid" name="userid" required pattern="[0-9a-zA-Z].{5,}" placeholder="Pick a username"><br>
                    <label for="signup-password">Password:</label><br>
                    <input type="password" id="signup-password" name="password" required
                            placeholder="Enter your password">
                     <img src="https://media.geeksforgeeks.org/wp-content/uploads/20210917145551/eye.png" 
                        width="5%" height="5%" style="display:inlinel margin-left: -1.5%; vertical-align: middle"
                        id="signup-enter">
                    <br> 
                    <p id="password-regex" style="color: red; display: none; font-size: 14px;">
                        Password need to be at least 1 character long with at least 1 number and a special character(!@#$%^/)
                    </p>
                    <label for="password-reenter">Re-enter your password:</label><br>
                    <input type="password" id="password-reenter" required
                            placeholder="Reenter your password">
                     <img src="https://media.geeksforgeeks.org/wp-content/uploads/20210917145551/eye.png" 
                        width="5%" height="5%" style="display:inlinel margin-left: -1.5%; vertical-align: middle"
                        id="signup-reenter">
                    <br>
        <p id="reenter-regex" style="color: red; display: none; font-size: 14px;">
                        Password need to be at least 10 character long with at least 1 number and a special character(!@#$%^/)
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
            `)
}
