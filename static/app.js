


const signupButton = document.getElementById("signup");
signupButton.addEventListener(`click`, function () {signup();});

const loginButton = document.getElementById("login")
loginButton.addEventListener(`click`, function (){login();});
    

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

            if (data.token){
                localStorage.setItem(`token`, data.token);
            }

        }catch(error){
            alert(`Error: ${error.message}`);
        }
        
    }

        