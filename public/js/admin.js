// public/js/admin.js
const ADMIN_API_BASE = '/api/v1/admin';

// Вызывается из auth.js после успешного ввода JWT
async function initAdminData() {
    loadAdminMovies();
    // Навешиваем слушатели на формы админки
    document.getElementById('add-movie-form').addEventListener('submit', addMovie);
    document.getElementById('add-show-form').addEventListener('submit', addShow);
}

async function loadAdminMovies() {
    const token = localStorage.getItem('admin_token');
    try {
        // Запрашиваем ВСЕ фильмы для админки (даже неактивные)
        const res = await fetch(`${ADMIN_API_BASE}/movies`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const movies = await res.json();
        
        const list = document.getElementById('admin-movies-list');
        const select = document.getElementById('show-movie-select');
        list.innerHTML = '';
        select.innerHTML = '';

        movies.forEach(movie => {
            // Заполняем селект для создания сеансов
            if(movie.is_active) {
                select.innerHTML += `<option value="${movie.id}">${movie.title}</option>`;
            }

            // Рендерим список с тумблерами скрытия
            const btnText = movie.is_active ? "Убрать из проката" : "Вернуть в прокат";
            const btnClass = movie.is_active ? "btn-danger" : "btn-success";
            
            list.innerHTML += `
                <div class="admin-movie-item">
                    <span><strong>${movie.title}</strong> ${movie.is_active ? '' : '[АРХИВ]'}</span>
                    <button class="${btnClass}" onclick="toggleMovie(${movie.id})">${btnText}</button>
                </div>`;
        });
    } catch (err) { console.error("Ошибка админки:", err); }
}

async function addMovie(e) {
    e.preventDefault();
    const token = localStorage.getItem('admin_token');
    const title = document.getElementById('m-title').value;
    const desc = document.getElementById('m-desc').value;
    const duration = parseInt(document.getElementById('m-duration').value);

    try {
        const res = await fetch(`${ADMIN_API_BASE}/movies`, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ title, description: desc, duration })
        });

        if (res.ok) {
            document.getElementById('add-movie-form').reset();
            loadAdminMovies();
        }
    } catch (err) { console.error(err); }
}

async function toggleMovie(id) {
    const token = localStorage.getItem('admin_token');
    try {
        await fetch(`${ADMIN_API_BASE}/movies/toggle`, {
            method: 'PUT',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ id: id })
        });
        loadAdminMovies();
    } catch (err) { console.error(err); }
}

async function addShow(e) {
    e.preventDefault();
    const token = localStorage.getItem('admin_token');
    const movieId = parseInt(document.getElementById('show-movie-select').value);
    const startTime = new Date(document.getElementById('show-time').value).toISOString();
    const price = parseFloat(document.getElementById('show-price').value);

    try {
        const res = await fetch(`${ADMIN_API_BASE}/shows`, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ movie_id: movieId, start_time: startTime, price })
        });

        if (res.ok) {
            alert("Сеанс успешно добавлен в расписание!");
            document.getElementById('add-show-form').reset();
        }
    } catch (err) { console.error(err); }
}