// models/user.go
package models

import "time"

type User struct {
    UserID    string    `json:"user_id"`
    Username  string    `json:"username" binding:"required"`
    Email     string    `json:"email" binding:"required"`
    Password  string    `json:"password" binding:"required"`
    Role      string    `json:"role" binding:"required"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
