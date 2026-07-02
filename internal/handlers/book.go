package handlers

import (
	"cinema_booking/db"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
)

type BookRequest struct {
	SeatIDs []int `json:"seat_ids"`
}

func BookSeats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if len(req.SeatIDs) == 0 {
		http.Error(w, "Некорректный список мест", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Print("Ошибка при начале транзакции: ", err)
		http.Error(w, "Ошибка при начале транзакции", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query :=
		`SELECT id, status 
	FROM seat 
	WHERE id = ANY($1) FOR UPDATE`

	rows, err := tx.QueryContext(ctx, query, pq.Array(req.SeatIDs))
	if err != nil {
		log.Print("Ошибка при выполнении запроса: ", err)
		http.Error(w, "Ошибка при выполнении запроса", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var id int
	var status string

	for rows.Next() {
		err := rows.Scan(&id, &status)
		if err != nil {
			log.Print("Ошибка при чтении статуса места: ", err)
			http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
			return
		}
		if status != "available" {
			http.Error(w, "Одно или несколько выбранных мест уже забронированы", http.StatusConflict)
			return
		}
	}
	rows.Close()

	updateQuery :=
		`UPDATE seat 
	SET status = 'booked' 
	WHERE id = ANY($1)`

	_, err = tx.ExecContext(ctx, updateQuery, pq.Array(req.SeatIDs))
	if err != nil {
		log.Print("Ошибка при обновлении статуса мест: ", err)
		http.Error(w, "Ошибка при бронировании мест", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Print("Ошибка при коммите транзакции: ", err)
		http.Error(w, "Ошибка при сохранении брони", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
