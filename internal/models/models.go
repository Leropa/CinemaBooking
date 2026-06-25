package models

import "time"

type User struct {
	ID           int    `json:"id"`
	UserName     string `json:"user_name"`
	PasswordHash string `json:"-"`
}

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	IsActive    bool   `json:"is_active"`
}

type Show struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movie_id"`
	StartTime time.Time `json:"start_time"`
	Price     float64   `json:"price"`
}

type Seat struct {
	ID     int    `json:"id"`
	ShowID int    `json:"show_id"`
	Row    int    `json:"row"`
	Number int    `json:"number"`
	Status string `json:"status"`
}

type Booking struct {
	ID            int       `json:"id"`
	ShowID        int       `json:"show_id"`
	SeatID        int       `json:"seat_id"`
	CustomerEmail string    `json:"customer_email"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
