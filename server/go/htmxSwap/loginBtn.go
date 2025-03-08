package htmxswap

import (
	"github.com/gin-gonic/gin"
)

func LoginButton(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, `
<form class="user-signin">
<div class="modal fade" id="staticBackdrop" data-bs-backdrop="static" data-bs-keyboard="true" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="false">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>
        <div class="body-footer">
          <div class="modal-body">
                    <h3>Log-In</h3>
                <label for="userid">Username:</label><br>
                <input type="text" id="userid" required placeholder="User name"><br>
                <label for="user-password">Password:</label><br>
                <input type="password" id="user-password" placeholder="Enter your password" required><br>
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
                hx-post="/user-log-in"
                hx-taget="header"
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
