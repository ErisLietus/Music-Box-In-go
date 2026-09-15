document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');

  if (token) {
    document.getElementById('auth-section').style.display = 'none';
    document.getElementById('after-login').style.display = 'block';
  } else {
    document.getElementById('auth-section').style.display = 'block';
    document.getElementById('after-login').style.display = 'none';
  }
});


const signupButton = document.getElementById("signup");
signupButton.addEventListener(`click`, function () {signup();});

const loginButton = document.getElementById("login")
loginButton.addEventListener(`click`, function (){login();});

const logoutButton = document.getElementById("logout");
logoutButton.addEventListener(`click`, function (){logout();});
    

async function signup() {
    console.log("This button was clicked WITH GUSTO!")
    const email = document.getElementById('email').value
    const password = document.getElementById('password').value

    try {
        const res = await fetch("/api/signup", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        });
        if (!res.ok){
            const data = await res.json();
            throw new Error(`Failed to create user: ${data.error}`);
            }
            console.log("User created!")
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
                }
        }catch(error){
            alert(`Error: ${error.message}`);
        }
        
    }

    function logout() {
        console.log("user has logged out")
        localStorage.removeItem('token');
        document.getElementById('auth-section').style.display = 'block';
         document.getElementById('after-login').style.display = 'none';
}

        