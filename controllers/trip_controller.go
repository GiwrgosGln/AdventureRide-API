package controllers

import (
	"adventureride/models"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TripController struct {
    Collection *mongo.Collection
}

func (tc *TripController) GetTrips(c *gin.Context) {
    trips := []models.Trip{}

    cursor, err := tc.Collection.Find(context.TODO(), bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trips"})
        return
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.Background()) {
        var trip models.Trip
        cursor.Decode(&trip)
        trips = append(trips, trip)
    }

    c.JSON(http.StatusOK, trips)
}

func (tc *TripController) GetTrip(c *gin.Context) {
    tripID := c.Param("trip_id")
    objectID, err := primitive.ObjectIDFromHex(tripID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID format"})
        return
    }
    
    var trip models.Trip
    err = tc.Collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&trip)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trip"})
        log.Println("Error fetching trip:", err)
        log.Println("Requested trip ID:", tripID)
        return
    }
    
    c.JSON(http.StatusOK, trip)
    
}



func (tc *TripController) CreateTrip(c *gin.Context) {
    var trip models.Trip
    if err := c.ShouldBindJSON(&trip); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // You can directly access the CreatorID, StartLocation, and EndLocation fields
    // from the 'trip' object as they are part of the JSON payload.

    _, err := tc.Collection.InsertOne(context.TODO(), trip)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trip"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Trip created"})
}



    func (tc *TripController) UpdateTrip(c *gin.Context) {
        tripID := c.Param("trip_id")
    
        // Convert the trip_id string to an ObjectID
        objectID, err := primitive.ObjectIDFromHex(tripID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID format"})
            return
        }
    
        var trip models.Trip
        if err := c.ShouldBindJSON(&trip); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        result, err := tc.Collection.ReplaceOne(context.TODO(), bson.M{"_id": objectID}, trip)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trip"})
            log.Println("Error updating trip:", err)
            log.Println("Requested trip ID:", tripID)
            log.Println("Request Body:", trip)
            return
        }
    
        // Log the result of the ReplaceOne operation
        log.Printf("ReplaceOne Result - Matched: %v, Modified: %v, UpsertedID: %v", result.MatchedCount, result.ModifiedCount, result.UpsertedID)
    
        c.JSON(http.StatusOK, gin.H{"message": "Trip updated"})
    }




    func (tc *TripController) DeleteTrip(c *gin.Context) {
        tripID := c.Param("trip_id")
        
        objectID, err := primitive.ObjectIDFromHex(tripID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID format"})
            log.Println("Error parsing trip ID:", err)
            return
        }
        
        _, err = tc.Collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trip"})
            log.Println("Error deleting trip:", err)
            log.Println("Requested trip ID:", tripID)
            return
        }
        
        c.JSON(http.StatusOK, gin.H{"message": "Trip deleted"})
    }
    
