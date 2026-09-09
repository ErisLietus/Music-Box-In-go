const signupButton = document.getElementById("signup");
signupButton.addEventListener(`click`, console.log("I hate this"));

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
        