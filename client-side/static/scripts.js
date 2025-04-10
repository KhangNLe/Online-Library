document.addEventListener('htmx:afterSwap', function(e) {
    if (e.detail.target.id === 'login-container' || e.detail.target.id === 'passwordChange') {
        const modal = new bootstrap.Modal(document.getElementById('staticBackdrop'));
        modal.show();
    }

    const loginBtn = document.getElementById("submit-login");
    if (loginBtn != null) {
        const pass = document.getElementById("user-password");
        const toggle = document.getElementById("login-password");
        toggle.addEventListener("click", (e) => {
            e.preventDefault();
            passwordVisibility(pass, toggle)
        });
    }

    const register = document.getElementById("submit-signup");
    if (register != null) {
        const password = document.getElementById("signup-password");
        const reenter = document.getElementById("password-reenter");
        const passToggle = document.getElementById("signup-enter");
        const reenterToggle = document.getElementById("signup-reenter");
        const errorMsg = document.getElementById("mismatch-pass");
        const passregex = document.getElementById("password-regex");
        const reenterRegex = document.getElementById("reenter-regex");
        const userName = document.getElementById("signup-userid");
        const userRegx = document.getElementById("userid-regex");

        passToggle.addEventListener("click", (e) => {
            e.preventDefault();
            passwordVisibility(password, passToggle);
        });
        reenterToggle.addEventListener("click", (e) => {
            e.preventDefault();
            passwordVisibility(reenter, reenterToggle);
        })

        userName.addEventListener("keyup", () => {
            regexChecking(userName, userRegx, 6);
        })

        password.addEventListener("keyup", () => {
            regexChecking(password, passregex, 10);
        });

        reenter.addEventListener("keyup", () => {
            regexChecking(reenter, reenterRegex, 10);

            if (password.value != reenter.value) {
                errorMsg.style.display = "block";
            } else {
                errorMsg.style.display = "none";
            }
        });
    }

    const searchBtn = document.getElementById("searchBtn")
    if (searchBtn != null) {
        searchBtn.addEventListener("submit", () => {
            const display = document.querySelector(".display")
            display.innerHTML = `
                <div class="spinner-border text-light" role="status">
                    <span class="visually-hidden">Loading...</span>
                </div>
            `
        })
    }

    const pageBtns = document.querySelectorAll(".pageBtn");
    if (pageBtns != null) {
        pageBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                window.scrollTo({
                    top: 0,
                    behavior: "smooth"
                })
            });
        });
    }

});

/**
 * @param {Element} text_box
 * @param {Element} visible_img 
 **/
function passwordVisibility(text_box, visible_img) {
    let type = (text_box.getAttribute("type") === "password") ? "text" : "password";
    text_box.setAttribute("type", type);

    if (visible_img.src.match("https://media.geeksforgeeks.org/wp-content/uploads/20210917145551/eye.png")) {
        visible_img.src = "https://media.geeksforgeeks.org/wp-content/uploads/20210917150049/eyeslash.png";
    } else {
        visible_img.src = "https://media.geeksforgeeks.org/wp-content/uploads/20210917145551/eye.png";
    }
}

/**
 * @param {Element} text_box
 * @param {Element} warning 
 **/
function regexChecking(text_box, warning, length) {
    let regx = createRegex(length);

    if (text_box.value == "") {
        warning.style.display = "none";
    } else if (!regx.test(text_box.value)) {
        warning.style.display = "block";
    } else {
        warning.style.display = "none";
    }
}

/**
 * @param {int} length
 * @returns RegExp
 **/
function createRegex(length) {
    return new RegExp(`^(?=.*\\d)(?=.*[a-zA-Z])(?=.*[!@#$%^/_]).{${length},}$`);
}

