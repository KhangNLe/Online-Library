package htmxswap

import "github.com/gin-gonic/gin"

func SignUpBtn(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, `
          <div class="modal-body">
                <h3>Sign-In</h3>
                    <label for="sign-up-userid">Username:</label><br>
                    <input type="text" id="sign-up-userid" required pattern="(?=.*[a-zA-Z0-9._-]).{5,}" placeholder="Pick a username"><br>
                    <label for="sign-up-password">Password:</label><br>
                    <input type="password" id="sign-up-password" 
                        pattern="(?=.*\d).{1,} (?=.*[a-zA-Z]).{8,}" placeholder="Enter your password"
                        required title="Password need to be at least 8 letters long with a number"><br> 
                    <label for="password-reenter">Re-enter your password:</label><br>
                    <input type="password" id="password-reenter"
                        pattern="(?=.*\d).{1,} (?=.*[a-zA-Z]).{8,}" placeholder="Enter your password"
                        required title="Password need to be at least 8 letters long with a number"><br> 
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
            <button type="submit" id="submit-signup" class="btn btn-primary">Register</button>
          </div>

            `)
}
