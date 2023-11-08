package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"adventureride/controllers"
	"adventureride/routes"
)

func main() {
    // Load environment variables from the .env file
    if err := godotenv.Load(); err != nil {
        fmt.Println("Error loading .env file")
        return
    }

    // Get the MongoDB URI from the environment variable
    mongoURI := os.Getenv("MONGODB_URI")

    // Initialize the Gin router
    r := gin.Default()

    // Connect to MongoDB Atlas
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
    if err != nil {
        fmt.Println("Failed to connect to MongoDB Atlas")
        return
    }
    defer client.Disconnect(context.TODO())

    // Create collections for users and trips
    userCollection := client.Database("AdventureRide").Collection("Users")
    tripCollection := client.Database("AdventureRide").Collection("Trips")

    // Create controllers for users and trips
    userController := &controllers.Controller{Collection: userCollection}
    tripController := &controllers.TripController{Collection: tripCollection}

    // Set up routes for User and Trip controllers
    routes.SetupUserRoutes(r, userController)
    routes.SetupTripRoutes(r, tripController)

    // Start the web server
    r.Run(":8080")
}
