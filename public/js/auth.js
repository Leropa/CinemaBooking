// public/js/auth.js
// Отвечает за хранение токена и перехват заголовков
function getToken() { return localStorage.getItem('admin_token'); }

document.addEventListener('DOMContentLoaded', () => {
    const token = getToken();
    if (token) {
        showDashboard();
    }

    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
});

async function handleLogin(e) {
    e.preventDefault();
    const user = document.getElementById('admin-username').value;
    const pass = document.getElementById('admin-password').value;

    try {
        const res = await fetch('/api/v1/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: user, password: pass })
        });

        if (res.ok) {
            const data = await res.json();
            localStorage.setItem('admin_token', data.token); // Сохраняем JWT
            showDashboard();
        } else {
            alert('Неверный логин или пароль персоналa!');
        }
    } catch (err) { console.error("Ошибка авторизации:", err); }
}

function showDashboard() {
    document.getElementById('auth-block').classList.add('hidden');
    document.getElementById('dashboard-block').classList.remove('hidden');
    // Сигнализируем скрипту админки загружать данные
    if (typeof initAdminData === 'function') initAdminData();
}