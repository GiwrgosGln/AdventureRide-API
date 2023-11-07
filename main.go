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

var collection *mongo.Collection

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
	collection = client.Database("AdventureRide").Collection("Users")

	// Create a controller and pass the collection
	ctrl := &controllers.Controller{Collection: collection}

	// Set up routes and pass the controller
	routes.SetupRoutes(r, ctrl)

	// Start the web server
	r.Run(":8080")
}
