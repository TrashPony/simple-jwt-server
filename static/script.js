const serverUrl = window.location.protocol + "//" + window.location.hostname + ':' + window.location.port

const api = {
    login: "/login",
    get_users: "/get_users"
}

function getToken() {
    fetch(serverUrl + api.login, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            user_name: document.getElementById("username").value,
            password: document.getElementById("password").value,
        })
    })
        .then(response => response.json())
        .then(result => {
            if (result.error) {
                console.error(result.error)
                return
            }

            if (result.token) {
                document.getElementById("token").innerText = result.token
                window.localStorage.setItem('x-auth-token', result.token)
            }
        })
}

function getUsers() {
    fetch(serverUrl + api.get_users, {
        headers: {
            Authorization: 'Bearer ' + window.localStorage.getItem('x-auth-token'),
        },
    })
        .then(response => response.json())
        .then(result => {
            if (result.error) {
                console.error(result.error)
                return
            }

            console.log(result)
        })
}