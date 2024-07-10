// models/home.go
package models

import "time"

type Home struct {
    HomeID      string    `json:"home_id"`
    OwnerID     string    `json:"owner_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Address     string    `json:"address"`
    City        string    `json:"city"`
    State       string    `json:"state"`
    ZipCode     string    `json:"zip_code"`
    Price       float64   `json:"price"`
    NumBedrooms int       `json:"num_bedrooms"`
    NumBathrooms int      `json:"num_bathrooms"`
    Sqft        int       `json:"sqft"`
    ImageURLs   []string  `json:"image_urls"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Status        string    `json:"status"` // new field for status: 'reserved' or 'available'
}