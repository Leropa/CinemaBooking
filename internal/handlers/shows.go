package handlers

import (
	"cinema_booking/db"
	"cinema_booking/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func GetShows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	movieID := r.URL.Query().Get("movie_id")

	query := `
		SELECT id, movie_id, start_time, price FROM show
		WHERE movie_id = $1
	`
	rows, err := db.DB.Query(query, movieID)
	if err != nil {
		http.Error(w, "Ошибка при выполнении запроса к базе данных", http.StatusInternalServerError)
		log.Print("Ошибка при выполнении запроса к базе данных", err)
		return
	}
	defer rows.Close()

	var s []models.Show
	for rows.Next() {
		var show models.Show

		err := rows.Scan(&show.ID, &show.MovieID, &show.StartTime, &show.Price)
		if err != nil {
			http.Error(w, "Ошибка при чтении данных из базы данных", http.StatusInternalServerError)
			log.Print("Ошибка при чтении данных из базы данных", err)
			return
		}

		s = append(s, show)
	}

	json.NewEncoder(w).Encode(s)
}
