document.addEventListener('htmx:afterSwap', function(e) {
    if (e.detail.target.id === 'login-container') {
        const modal = new bootstrap.Modal(document.getElementById('staticBackdrop'));
        modal.show();
    }
    const loginBtn = document.getElementById("submit-login");
    if (loginBtn != null) {
        loginBtn.addEventListener('click', function(e) {
            e.preventDefault();
            const userId = document.getElementById('userid');
            const pass = document.getElementById('user-password');
            console.log(userId.value);
            console.log(pass.value);
            fetch("http://localhost:6969/user-log-in", {
                method: "POST", headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    userName: userId.value,
                    password: pass.value,
                })
            })

                .then(respone => respone.json())
                .then(data => console.log(data))
                .catch(error => console.error("Error: ", error));
        });
    }
});


