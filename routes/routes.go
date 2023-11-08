package routes

import (
	"adventureride/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, userCtrl *controllers.Controller) {
    userGroup := r.Group("/users")
    {
        userGroup.POST("/register", userCtrl.RegisterHandler)
        userGroup.POST("/login", userCtrl.LoginHandler)
    }
}

func SetupTripRoutes(r *gin.Engine, tripCtrl *controllers.TripController) {
    tripGroup := r.Group("/trips")
    {
        tripGroup.GET("/", tripCtrl.GetTrips)
        tripGroup.GET("/:trip_id", tripCtrl.GetTrip)
        tripGroup.POST("/", tripCtrl.CreateTrip)
        tripGroup.PUT("/:trip_id", tripCtrl.UpdateTrip)
        tripGroup.DELETE("/:trip_id", tripCtrl.DeleteTrip)
    }
}
