// controllers/reservation.go
package controllers

import (
	"context"
	"homefinder/db"
	"homefinder/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()
	reservation.Status = "pending"

	// Check the status of the home
	var homeStatus string
	err := db.Conn.QueryRow(context.Background(), "SELECT status FROM homes WHERE home_id = $1", reservation.HomeID).Scan(&homeStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching home status"})
		return
	}

	// Only allow reservation if the home status is 'available'
	if homeStatus != "available" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create reservation. Home is not available."})
		return
	}

	sql := `INSERT INTO reservations (home_id, user_id, start_date, end_date, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING reservation_id`
	err = db.Conn.QueryRow(context.Background(), sql, reservation.HomeID, reservation.UserID, reservation.StartDate, reservation.EndDate, reservation.Status, reservation.CreatedAt, reservation.UpdatedAt).Scan(&reservation.ReservationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating reservation"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}



func ApproveReservation(c *gin.Context) {
	// Extract home ID and reservation ID from request parameters
	homeID := c.Param("home_id")
	reservationID := c.Param("reservation_id")

	// Check if the current user is the owner of the home
	userID := c.MustGet("user_id").(string)
	var ownerID string
	err := db.Conn.QueryRow(context.Background(), "SELECT owner_id FROM homes WHERE home_id = $1", homeID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching home owner"})
		return
	}

	if userID != ownerID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to approve reservations for this home"})
		return
	}

	// Update the status of the home to 'reserved'
	_, err = db.Conn.Exec(context.Background(), "UPDATE homes SET status = 'reserved' WHERE home_id = $1", homeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating home status"})
		return
	}

	// Update the status of the reservation to 'confirmed'
	_, err = db.Conn.Exec(context.Background(), "UPDATE reservations SET status = 'confirmed' WHERE reservation_id = $1", reservationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating reservation status"})
		return
	}

	// Optionally, you may want to send a notification to the renter confirming the reservation approval

	c.JSON(http.StatusOK, gin.H{"message": "Reservation approved successfully"})
}