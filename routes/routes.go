// // routes/routes.go
// package routes

// import (
//     "homefinder/controllers"
//     "github.com/gin-gonic/gin"
// )

// func InitializeRoutes(router *gin.Engine) {
//     router.POST("/register", controllers.Register)
// 	router.POST("/reservations", controllers.CreateReservation)
//     // Add more routes here
// }

package routes

import (
    "homefinder/controllers"
    "homefinder/middleware"

    "github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
    router.POST("/register", controllers.Register)
    router.POST("/login", controllers.Login)

    // Public routes for fetching homes
    router.GET("/homes", controllers.GetHomes)
    router.GET("/homes/:home_id", controllers.GetHomeByID)

    // Authenticated routes for home management
    auth := router.Group("/")
    auth.Use(middleware.JWTAuthMiddleware())
    {
        auth.POST("/homes", middleware.OwnerRoleMiddleware(), controllers.CreateHome)
        auth.PUT("/homes/:home_id", middleware.OwnerRoleMiddleware(), controllers.UpdateHome)
        auth.DELETE("/homes/:home_id", middleware.OwnerRoleMiddleware(), controllers.DeleteHome)
        auth.POST("/reservations", controllers.CreateReservation)
		auth.PUT("/reservations/:home_id/:reservation_id/approve", middleware.OwnerRoleMiddleware(), controllers.ApproveReservation)
    }
}

