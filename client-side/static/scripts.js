const loginBtn = document.getElementById("log-in");
const homeBtn = document.getElementById('home');
const seachBtn = document.getElementById('search');
const booksBtn = document.getElementById('mybook');
const recBtn = document.getElementById('recommend');
const aboutBtn = document.getElementById('about');
const loginContainer = document.querySelector('.login-btn');

document.addEventListener('DOMContentLoaded', function() {
    loginContainer.innerHTML = `
<button type="button" class="btn btn-dark" data-bs-target="#staticBackdrop" data-bs-toggle="modal" id="log-in">Log-In</button>

<form class="loggin">
<div class="modal fade" id="staticBackdrop" data-bs-backdrop="static" data-bs-keyboard="true" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="false">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>
      <div class="modal-body">
                <h3>Log-In</h3>
            <label for="userid">Username:</label><br>
            <input type="text" id="userid" required placeholder="User name"><br>
            <label for="password">Password:</label><br>
            <input type="password" id="password" placeholder="Enter your password" required><br>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        <button type="button" class="btn btn-info" id="sign-up">Sign Up</button>
        <button type="submit" id="submit-login" class="btn btn-primary">Log-in</button>
      </div>
    </div>
  </div>
</div>
</form>
    `;
    userSignUp();

});


function userSignUp() {
    let signUpBtn = document.getElementById('sign-up');
    let body = document.querySelector(".modal-body");
    let footer = document.querySelector('.modal-footer');
    signUpBtn.addEventListener('click', () => {
        body.innerHTML = `
        <h3>Sign-In</h3>
        <label for="sign-up-userid">Username:</label><br>
        <input type="text" id="sign-up-userid" required pattern="(?=.*[a-zA-Z0-9._-]).{5,}" placeholder="Pick a username"><br>
        <label for="sign-up-password">Password:</label><br>
        <input type="password" id="sign-up-password" pattern="(?=.*\d)(?=.*[a-zA-Z]).{8,}" placeholder="Enter your password" required><br> 
        <label for="password-reenter">Re-enter your password:</label><br>
        <input type="password" id="password-reenter" pattern="(?=.*\d)(?=.*[a-zA-Z]).{8,}" placeholder="Enter your password" required><br> 
        `;
        footer.innerHTML = `
                <button type="button" id="closebtn" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
            <button type="button" class="btn btn-primary" id="sign-up">Register</button>
        `;
        let closeBtn = document.getElementById('closebtn');
        closeBtn.addEventListener('click', () => {
            body.innerHTML = `
                 <h3>Log-In</h3>
            <label for="userid">Username:</label><br>
            <input type="text" id="userid" required placeholder="User name"><br>
            <label for="password">Password:</label><br>
            <input type="password" id="password" placeholder="Enter your password" required><br>
            `;

            footer.innerHTML = `
               <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                <button type="button" class="btn btn-info" id="sign-up">Sign Up</button>
                <button type="submit" id="submit-login" class="btn btn-primary">Log-in</button>
            `;
        });
    });
}
