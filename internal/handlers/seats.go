package handlers

import (
	"cinema_booking/db"
	"cinema_booking/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func GetSeats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	showID := r.URL.Query().Get("show_id")

	query := `
		SELECT id, show_id, row, number, status FROM seat
		WHERE show_id = $1
		ORDER BY row, number
	`

	rows, err := db.DB.QueryContext(r.Context(), query, showID)
	if err != nil {
		http.Error(w, "Ошибка при выполнении запроса к базе данных", http.StatusInternalServerError)
		log.Print("Ошибка при выполнении запроса к базе данных: ", err)
		return
	}

	var s []models.Seat
	for rows.Next() {
		var seat models.Seat
		err := rows.Scan(&seat.ID, &seat.ShowID, &seat.Row, &seat.Number, &seat.Status)
		if err != nil {
			rows.Close()
			http.Error(w, "Ошибка при чтении данных из базы данных", http.StatusInternalServerError)
			log.Print("Ошибка при чтении данных из базы данных: ", err)
			return
		}
		s = append(s, seat)
	}
	rows.Close()

	if len(s) == 0 {
		for row := 1; row <= 5; row++ {
			for number := 1; number <= 10; number++ {
				insertQuery :=
					`INSERT INTO seat (show_id, row, number, status) 
				VALUES ($1, $2, $3, 'available')`
				_, err := db.DB.ExecContext(r.Context(), insertQuery, showID, row, number)
				if err != nil {
					http.Error(w, "Ошибка при генерации мест в зале", http.StatusInternalServerError)
					log.Print("Ошибка при вставке места: ", err)
					return
				}
			}
		}

		finalRows, err := db.DB.QueryContext(r.Context(), query, showID)
		if err != nil {
			http.Error(w, "Ошибка при получении сгенерированных мест", http.StatusInternalServerError)
			log.Print("Ошибка после генерации мест: ", err)
			return
		}
		defer finalRows.Close()

		for finalRows.Next() {
			var seat models.Seat
			if err := finalRows.Scan(&seat.ID, &seat.ShowID, &seat.Row, &seat.Number, &seat.Status); err != nil {
				http.Error(w, "Ошибка при чтении сгенерированных мест", http.StatusInternalServerError)
				return
			}
			s = append(s, seat)
		}
	}

	json.NewEncoder(w).Encode(s)
}
