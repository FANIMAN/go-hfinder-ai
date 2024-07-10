// models/reservation.go
package models

import "time"

type Reservation struct {
	ReservationID string    `json:"reservation_id"`
	HomeID        string    `json:"home_id"`
	UserID        string    `json:"user_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
