// main.go
package main

import (
	"homefinder/db"
	"homefinder/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	
    // Initialize the database connection
    db.ConnectDB()
    defer db.CloseDB()

    // Set up the Gin router
    r := gin.Default()
    routes.InitializeRoutes(r)

    // Graceful shutdown
    go func() {
        if err := r.Run(":8080"); err != nil {
            log.Fatalf("Unable to start server: %v\n", err)
        }
    }()

    // Wait for interrupt signal to gracefully shut down the server
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")
}
