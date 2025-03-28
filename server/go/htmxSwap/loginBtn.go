package htmxswap

import (
	"github.com/gin-gonic/gin"
)

func LoginButton(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, `
<form class="user-signin" 
        hx-post="/user-login" 
        hx-on::after-request=" 
            if (event.detail.xhr.status >= 400) { document.getElementById('login-warning').innerHTML = event.detail.xhr.responseText; }
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
                    <h3>Log-In</h3>
                <label for="userid">Username:</label><br>
                <input type="text" id="userid"
                        name="userid" required placeholder="User name"><br>
                <label for="user-password">Password:</label><br>
                <input type="password" id="user-password"
                        name="password" placeholder="Enter your password" required>
                <img src="https://media.geeksforgeeks.org/wp-content/uploads/20210917145551/eye.png" 
                        width="5%" height="5%" style="display:inlinel margin-left: -1.5%; vertical-align: middle"
                        id="login-password">
                <br>
            <div id="login-warning">
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
            <button
                type="button"
                class="btn btn-info"
                hx-get="/sign-up"
                hx-target=".body-footer"
                hx-swap="innerHTML"
                hx-trigger="click"
                id="sign-up"
                >Sign Up
            </button>
            <button type="submit" 
                id="submit-login" 
                class="btn btn-primary"
                >Log-in
            </button>
        </div>
    </div>
    </div>
  </div>
</div>
</form>
            `)
}
