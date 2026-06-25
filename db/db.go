package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error
	infoDB := "host=localhost port=5432 user=postgres password=leropa dbname=cinema_booking sslmode=disable"

	DB, err = sql.Open("postgres", infoDB)
	if err != nil {
		fmt.Print("Не получилось открыть Базу дынных:", err)
		return
	}

	query := `CREATE TABLE IF NOT EXISTS admin_users(
			id SERIAL PRIMARY KEY,
			user_name VARCHAR(25),
			password_hash VARCHAR(15)
		);

		CREATE TABLE IF NOT EXISTS movie(
			id SERIAL PRIMARY KEY,
			title VARCHAR(50),
			description VARCHAR,
			duration INT,
			is_active BOOL
		);

		CREATE TABLE IF NOT EXISTS show(
			id SERIAL PRIMARY KEY,
			movie_id INT REFERENCES movie(id) ON DELETE RESTRICT,
			start_time TIMESTAMP,
			price NUMERIC(10, 2) DEFAULT 0.00
		);

		CREATE TABLE IF NOT EXISTS seat(
			id SERIAL PRIMARY KEY,
			show_id INT REFERENCES show(id) ON DELETE RESTRICT,
			row INT,
			number INT,
			status VARCHAR(10)
		);
		
		CREATE TABLE IF NOT EXISTS booking(
			id SERIAL PRIMARY KEY,
			show_id INT REFERENCES show(id) ON DELETE RESTRICT,
			seat_id INT REFERENCES seat(id) ON DELETE RESTRICT,
			customer_email VARCHAR(50),
			status VARCHAR(10),
			created_at TIMESTAMP DEFAULT NOW()
		);
	`

	_, err = DB.Exec(query)
	if err != nil {
		log.Print("Не удалось про иницилизировать таблицы:", err)
	}

}
