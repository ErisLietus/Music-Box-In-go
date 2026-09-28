document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');

  if (token) {
    const isValid = await checktoken(token);
    if (isValid) {
    document.getElementById('auth-section').style.display = 'none';
    document.getElementById('after-login').style.display = 'block';
    } else {
        localStorage.removeItem(`token`)
        document.getElementById('login-form').style.display = 'block';
        document.getElementById("signup-form").style.display = 'none';
        document.getElementById('after-login').style.display = 'none';
        document.getElementById(`create-media`).style.display = 'none';
    }
  } else {
     document.getElementById('login-form').style.display = 'block';
        document.getElementById("signup-form").style.display = 'none';
        document.getElementById('after-login').style.display = 'none';
        document.getElementById(`create-media`).style.display = 'none';
  }
});


const beforeSignupButton = document.getElementById("before-sign-up")
beforeSignupButton.addEventListener(`click`, function() {
    document.getElementById('login-form').style.display = 'none';
    document.getElementById("signup-form").style.display = 'block';
})

document.getElementById("back-to-login").addEventListener("click", () => {
    document.getElementById('signup-form').style.display = 'none';
    document.getElementById('login-form').style.display = 'block';
});

const signupButton = document.getElementById("signup");
signupButton.addEventListener(`click`, function () {signup();});

const loginButton = document.getElementById("login")
loginButton.addEventListener(`click`, function (){login();});

const logoutButton = document.getElementById("logout");
logoutButton.addEventListener(`click`, function (){logout();});

const createPlaylistButton = document.getElementById("create-playlist-button");
createPlaylistButton.addEventListener(`click`, function(){createPlaylist();});



async function checktoken(token) {
    const res = await fetch('/api/checkUser', {
        headers: {Authorization: `Bearer ${token}`}
    });
    return res.ok
}

async function signup() {
    console.log("This button was clicked WITH GUSTO!")
    const email = document.getElementById('signup-email').value
    const password = document.getElementById('signup-password').value
    const username = document.getElementById(`signup-username`).value

    try {
        const res = await fetch("/api/signup", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password , username}),
        });
        if (!res.ok){
            const data = await res.json();
            throw new Error(`Failed to create user: ${data.error}`);
            }
            console.log("User created!")
            await login();
        }catch (error){
             alert(`Error: ${error.message}`);
        }
    }

    async function login() {
        console.log("Login js was hit")
        const email = document.getElementById(`email`).value;
        const password = document.getElementById(`password`).value;

        try {
            const res = await fetch(`/api/login`, {
                method: `POST`,
                headers: {
                    'Content-Type': `application/json`,
                },
                body: JSON.stringify({email, password}),
            });
            const data = await res.json();
            if (!res.ok) {
                throw new Error(`Failed to login: ${data.error}`);
            
            }
                if (data.token) {
                    console.log("User has logged out")
                    localStorage.setItem('token', data.token);
                    document.getElementById('auth-section').style.display = 'none';
                    document.getElementById('after-login').style.display = 'block';
                    await getPlaylists()
                } else {
                    alert('Login failed. Please check your credentials.');
                }
        }catch(error){
            alert(`Error: ${error.message}`);
        }
        
    }

    function logout() {
        console.log("user has logged out")
        localStorage.removeItem('token');
        document.getElementById('login-form').style.display = 'block';
        document.getElementById('signup-form').style.display = 'none';
        document.getElementById('after-login').style.display = 'none';
}

    async function createPlaylist() {
        const name = document.getElementById("name").value
        const isPublic = document.getElementById("public").checked
        const allowCollabEdits = document.getElementById("collab").checked

        try {
            const res = await fetch("/api/playlists", {
                method: "POST",
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${localStorage.getItem(`token`)}`, 
                }, 
                body: JSON.stringify({name, isPublic, allowCollabEdits}),
            });
            if (!res.ok) {
                throw new Error("Failed to create playlist")
            }
            console.log("It uploaded")
           document.getElementById("yay-button").style.display = `block`; 
           document.getElementById("create-media").style.display = `block`;
            
        }catch(error){
            alert(`Error: ${error.message}`);
            }
        }

        async function getPlaylists() {
            try {
                const res = await fetch("/api/getUserPlaylists", {
                    method: 'GET',
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('token')}`,
                     }
                    });
                        if (!res.ok){
                            const data = await res.json();
                            throw new error(`Failed to get playlists, Error: ${data.error}`);
                        }
                        const playlists = await res.json();
                        const playlistList = document.getElementById("playlist-list");
                        playlistList.innerHTML = ""
                        for (const playlist of playlists){
                            console.log(playlist.Name)
                            const listItem = document.createElement("li");
                            listItem.textContent = playlist.Name;
                            playlistList.appendChild(listItem)
                        }
                    }catch (error){
                        alert(`Error json: ${error.message}`);
                    }
        }
            
        

    
        