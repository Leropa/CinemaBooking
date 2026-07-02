package handlers

import (
	"cinema_booking/db"
	"cinema_booking/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT id, title, description, duration, is_active 
	FROM movie 
	WHERE is_active = true`

	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Ошибка при выполнении запроса к базе данных", http.StatusInternalServerError)
		log.Print("Ошибка при выполнении запроса к базе данных", err)
		return
	}
	defer rows.Close()

	var m []models.Movie

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.IsActive)
		if err != nil {
			http.Error(w, "Ошибка при чтении данных из базы данных", http.StatusInternalServerError)
			log.Print("Ошибка при чтении данных из базы данных", err)
			return
		}
		m = append(m, movie)
	}

	json.NewEncoder(w).Encode(m)
}
