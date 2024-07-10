// controllers/user.go
package controllers

import (
	"context"
	"homefinder/db"
	"homefinder/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Role != "renter" && user.Role != "homeowner" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role value"})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		return
	}

	log.Printf("Raw password during registration: %s", user.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	log.Printf("Hashed password during registration: %s", user.Password)

	sql := `INSERT INTO users (username, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id`
	err = db.Conn.QueryRow(context.Background(), sql, user.Username, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.UserID)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user", "details": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	log.Printf("Env variable JWT_SECRET: %s", os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error generating token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Login(c *gin.Context) {
	var loginDetails struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Raw password provided during login: %s", loginDetails.Password)

	var user models.User
	sql := `SELECT user_id, password FROM users WHERE email = $1`
	err := db.Conn.QueryRow(context.Background(), sql, loginDetails.Email).Scan(&user.UserID, &user.Password)
	if err != nil {
		log.Printf("Error on login1: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	log.Printf("Hashed password from DB: %s", user.Password)
	log.Printf("Password provided during login: %s", loginDetails.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDetails.Password)); err != nil {
		log.Printf("Error on login2: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	log.Printf("Env variable JWT_SECRET: %s", os.Getenv("JWT_SECRET"))

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error generating token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
