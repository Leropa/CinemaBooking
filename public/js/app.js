// public/js/app.js
const API_BASE = '/api/v1';
let selectedShowId = null;
let selectedSeatId = null;
let countdownInterval = null;

document.addEventListener('DOMContentLoaded', () => {
    loadMovies();
    
    // Кнопки возврата
    document.getElementById('back-to-movies').addEventListener('click', () => switchSection('movies-section'));
    document.getElementById('back-to-shows').addEventListener('click', () => switchSection('shows-section'));
    
    // Форма покупки
    document.getElementById('purchase-form').addEventListener('submit', handlePurchase);
});

function switchSection(sectionId) {
    ['movies-section', 'shows-section', 'hall-section'].forEach(id => {
        document.getElementById(id).classList.add('hidden');
    });
    document.getElementById(sectionId).classList.remove('hidden');
}

// 1. Загрузка фильмов
async function loadMovies() {
    try {
        const res = await fetch(`${API_BASE}/movies`);
        const movies = await res.json();
        const grid = document.getElementById('movies-grid');
        grid.innerHTML = '';
        
        movies.forEach(movie => {
            grid.innerHTML += `
                <div class="movie-card">
                    <div>
                        <h3>${movie.title}</h3>
                        <p>${movie.description}</p>
                    </div>
                    <small style="display:block; margin-bottom:10px;">Длительность: ${movie.duration} мин.</small>
                    <button class="btn-primary" onclick="showMovieSchedule(${movie.id}, '${movie.title}')">Выбрать сеанс</button>
                </div>`;
        });
    } catch (err) { console.error("Ошибка загрузки фильмов:", err); }
}

// 2. Загрузка сеансов фильма
async function showMovieSchedule(movieId, title) {
    try {
        document.getElementById('selected-movie-title').innerText = title;
        const res = await fetch(`${API_BASE}/shows?movie_id=${movieId}`);
        const shows = await res.json();
        const list = document.getElementById('shows-list');
        list.innerHTML = '';

        if(shows.length === 0) { list.innerHTML = '<p>Сеансов на этот фильм пока нет.</p>'; }
        
        shows.forEach(show => {
            const time = new Date(show.start_time).toLocaleString('ru-RU', { hour: '2-digit', minute: '2-digit', day: 'numeric', month: 'short' });
            list.innerHTML += `
                <div class="show-item">
                    <div><strong>Время: ${time}</strong></div>
                    <div>Цена: ${show.price} руб.</div>
                    <button class="btn-primary" onclick="openHallMap(${show.id}, ${show.price})">Выбрать место</button>
                </div>`;
        });
        switchSection('shows-section');
    } catch (err) { console.error("Ошибка загрузки сеансов:", err); }
}

// 3. Загрузка интерактивного зала (Синхронизировано со статусом 'available' из Go)
async function openHallMap(showId, price) {
    selectedShowId = showId;
    document.getElementById('booking-form-container').classList.add('hidden');
    if(countdownInterval) clearInterval(countdownInterval);

    try {
        const res = await fetch(`${API_BASE}/seats?show_id=${showId}`);
        const seats = await res.json();
        const grid = document.getElementById('hall-grid');
        grid.innerHTML = '';

        seats.forEach(seat => {
            const seatEl = document.createElement('div');
            seatEl.className = `seat ${seat.status}`;
            seatEl.innerText = seat.number;
            
            // Проверяем статус 'available', который генерирует наш бэкенд на Go
            if (seat.status === 'available') {
                seatEl.onclick = () => selectSeat(seatEl, seat.id, seat.row, seat.number, price);
            }
            grid.appendChild(seatEl);
        });
        switchSection('hall-section');
    } catch (err) { console.error("Ошибка карты зала:", err); }
}

// 4. Клик по свободному месту
async function selectSeat(element, seatId, row, number, price) {
    // Снимаем выделение с прошлого выбранного места локально
    document.querySelectorAll('.seat.selected').forEach(el => el.classList.remove('selected'));
    element.classList.add('selected');
    
    selectedSeatId = seatId;
    document.getElementById('selected-seat-info').innerText = `Ряд ${row}, Место ${number}`;
    document.getElementById('ticket-price').innerText = price;
    document.getElementById('booking-form-container').classList.remove('hidden');

    // Запускаем локальный таймер на 10 минут
    startCountdown(600);
}

function startCountdown(duration) {
    if(countdownInterval) clearInterval(countdownInterval);
    let timer = duration, minutes, seconds;
    const display = document.getElementById('timer-count');
    
    countdownInterval = setInterval(() => {
        minutes = parseInt(timer / 60, 10);
        seconds = parseInt(timer % 60, 10);
        minutes = minutes < 10 ? "0" + minutes : minutes;
        seconds = seconds < 10 ? "0" + seconds : seconds;
        display.textContent = minutes + ":" + seconds;

        if (--timer < 0) {
            clearInterval(countdownInterval);
            alert("Время бронирования истекло! Место освобождено.");
            openHallMap(selectedShowId, document.getElementById('ticket-price').innerText);
        }
    }, 1000);
}

// 5. Покупка (Финальное подтверждение на защищенный транзакцией Go-бэк)
async function handlePurchase(e) {
    e.preventDefault();

    try {
        // Делаем запрос на наш эндпоинт бронирования
        const res = await fetch(`${API_BASE}/book`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                // Передаем ID выбранного сиденья массивом, как требует BookRequest в Go
                seat_ids: [selectedSeatId] 
            })
        });

        if (res.ok) {
            clearInterval(countdownInterval);
            alert("Билет успешно куплен!");
            switchSection('movies-section');
            loadMovies();
        } else {
            // Сюда попадем, если бэкенд выкинул ошибку (например, статус 409 Conflict при гонке)
            alert("Ошибка: Это место уже забронировано другим пользователем.");
        }
    } catch (err) { console.error("Ошибка покупки:", err); }
}