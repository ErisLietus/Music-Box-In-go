
const fetch = (await import('node-fetch')).default


document.addEventListener('DOMContentLoaded', onclick)



document
  .getElementById("login-form")
  .addEventListener("submit", async (event) => {
    event.preventDefault();
    await login();
  });




async function login() {
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`Failed to login: ${data.error}`);
    }

    if (data.token) {
      localStorage.setItem('token', data.token);
      document.getElementById('auth-section').style.display = 'none';
      document.getElementById('video-section').style.display = 'block';
    } else {
      alert('Login failed. Please check your credentials.');
    }
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

async function signup() {
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  try {
    const res = await fetch('/api/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });
    if (!res.ok) {
      const data = await res.json();
      throw new Error(`Failed to create user: ${data.error}`);
    }
    console.log('User created!');
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}