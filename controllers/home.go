package controllers

import (
	"context"
	"homefinder/db"
	"homefinder/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// func CreateHome(c *gin.Context) {
// 	var home models.Home
// 	if err := c.ShouldBindJSON(&home); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	userID := c.MustGet("user_id").(string)
// 	home.OwnerID = userID
// 	home.CreatedAt = time.Now()
// 	home.UpdatedAt = time.Now()

// 	sql := `INSERT INTO homes (owner_id, title, description, address, city, state, zip_code, price, num_bedrooms, num_bathrooms, sqft, image_urls, created_at, updated_at)
//             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING home_id`
// 	err := db.Conn.QueryRow(context.Background(), sql, home.OwnerID, home.Title, home.Description, home.Address, home.City, home.State, home.ZipCode, home.Price, home.NumBedrooms, home.NumBathrooms, home.Sqft, home.ImageURLs, home.CreatedAt, home.UpdatedAt).Scan(&home.HomeID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating home"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, home)
// }

func CreateHome(c *gin.Context) {
	var home models.Home
	if err := c.ShouldBindJSON(&home); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(string)
	home.OwnerID = userID
	home.CreatedAt = time.Now()
	home.UpdatedAt = time.Now()
	home.Status = "available" // Set default status to 'available'

	sql := `INSERT INTO homes (owner_id, title, description, address, city, state, zip_code, price, num_bedrooms, num_bathrooms, sqft, image_urls, created_at, updated_at, status) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING home_id`
	err := db.Conn.QueryRow(context.Background(), sql, home.OwnerID, home.Title, home.Description, home.Address, home.City, home.State, home.ZipCode, home.Price, home.NumBedrooms, home.NumBathrooms, home.Sqft, home.ImageURLs, home.CreatedAt, home.UpdatedAt, home.Status).Scan(&home.HomeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating home"})
		return
	}

	c.JSON(http.StatusOK, home)
}

func GetHomes(c *gin.Context) {
	rows, err := db.Conn.Query(context.Background(), "SELECT * FROM homes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching homes"})
		return
	}
	defer rows.Close()

	homes := []models.Home{}
	for rows.Next() {
		var home models.Home
		err := rows.Scan(&home.HomeID, &home.OwnerID, &home.Title, &home.Description, &home.Address, &home.City, &home.State, &home.ZipCode, &home.Price, &home.NumBedrooms, &home.NumBathrooms, &home.Sqft, &home.ImageURLs, &home.CreatedAt, &home.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning home"})
			return
		}
		homes = append(homes, home)
	}

	c.JSON(http.StatusOK, homes)
}

func GetHomeByID(c *gin.Context) {
	homeID := c.Param("home_id")

	var home models.Home
	sql := `SELECT * FROM homes WHERE home_id = $1`
	err := db.Conn.QueryRow(context.Background(), sql, homeID).Scan(&home.HomeID, &home.OwnerID, &home.Title, &home.Description, &home.Address, &home.City, &home.State, &home.ZipCode, &home.Price, &home.NumBedrooms, &home.NumBathrooms, &home.Sqft, &home.ImageURLs, &home.CreatedAt, &home.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Home not found"})
		return
	}

	c.JSON(http.StatusOK, home)
}

// func UpdateHome(c *gin.Context) {
//     var home models.Home
//     if err := c.ShouldBindJSON(&home); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     homeID := c.Param("home_id")
//     home.UpdatedAt = time.Now()

//     sql := `UPDATE homes SET title=$1, description=$2, address=$3, city=$4, state=$5, zip_code=$6, price=$7, num_bedrooms=$8, num_bathrooms=$9, sqft=$10, image_urls=$11, updated_at=$12 WHERE home_id=$13`
//     _, err := db.Conn.Exec(context.Background(), sql, home.Title, home.Description, home.Address, home.City, home.State, home.ZipCode, home.Price, home.NumBedrooms, home.NumBathrooms, home.Sqft, home.ImageURLs, home.UpdatedAt, homeID)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating home"})
//         return
//     }

//     c.JSON(http.StatusOK, home)
// }

func UpdateHome(c *gin.Context) {
	homeID := c.Param("home_id")

	// Fetch the existing home details
	var existingHome models.Home
	sql := `SELECT owner_id, title, description, address, city, state, zip_code, price, num_bedrooms, num_bathrooms, sqft, image_urls, created_at, updated_at FROM homes WHERE home_id=$1`
	err := db.Conn.QueryRow(context.Background(), sql, homeID).Scan(&existingHome.OwnerID, &existingHome.Title, &existingHome.Description, &existingHome.Address, &existingHome.City, &existingHome.State, &existingHome.ZipCode, &existingHome.Price, &existingHome.NumBedrooms, &existingHome.NumBathrooms, &existingHome.Sqft, &existingHome.ImageURLs, &existingHome.CreatedAt, &existingHome.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching home details"})
		return
	}

	// Bind the JSON input
	var inputHome models.Home
	if err := c.ShouldBindJSON(&inputHome); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the provided fields
	if inputHome.Title != "" {
		existingHome.Title = inputHome.Title
	}
	if inputHome.Description != "" {
		existingHome.Description = inputHome.Description
	}
	if inputHome.Address != "" {
		existingHome.Address = inputHome.Address
	}
	if inputHome.City != "" {
		existingHome.City = inputHome.City
	}
	if inputHome.State != "" {
		existingHome.State = inputHome.State
	}
	if inputHome.ZipCode != "" {
		existingHome.ZipCode = inputHome.ZipCode
	}
	if inputHome.Price != 0 {
		existingHome.Price = inputHome.Price
	}
	if inputHome.NumBedrooms != 0 {
		existingHome.NumBedrooms = inputHome.NumBedrooms
	}
	if inputHome.NumBathrooms != 0 {
		existingHome.NumBathrooms = inputHome.NumBathrooms
	}
	if inputHome.Sqft != 0 {
		existingHome.Sqft = inputHome.Sqft
	}
	if inputHome.ImageURLs != nil {
		existingHome.ImageURLs = inputHome.ImageURLs
	}

	existingHome.UpdatedAt = time.Now()

	// Execute the update query
	sql = `UPDATE homes SET title=$1, description=$2, address=$3, city=$4, state=$5, zip_code=$6, price=$7, num_bedrooms=$8, num_bathrooms=$9, sqft=$10, image_urls=$11, updated_at=$12 WHERE home_id=$13`
	_, err = db.Conn.Exec(context.Background(), sql, existingHome.Title, existingHome.Description, existingHome.Address, existingHome.City, existingHome.State, existingHome.ZipCode, existingHome.Price, existingHome.NumBedrooms, existingHome.NumBathrooms, existingHome.Sqft, existingHome.ImageURLs, existingHome.UpdatedAt, homeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating home"})
		return
	}

	// Set the home_id explicitly in the response
	existingHome.HomeID = homeID

	c.JSON(http.StatusOK, existingHome)
}

func DeleteHome(c *gin.Context) {
	homeID := c.Param("home_id")

	sql := `DELETE FROM homes WHERE home_id=$1`
	_, err := db.Conn.Exec(context.Background(), sql, homeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting home"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Home deleted successfully"})
}
