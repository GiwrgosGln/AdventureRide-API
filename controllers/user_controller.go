package controllers

import (
	"context"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"adventureride/models"
	"adventureride/utils"

	"github.com/dgrijalva/jwt-go"
)

// Define a secret key for signing JWT tokens.
var jwtSecret = []byte("your-secret-key")

type Controller struct {
	Collection *mongo.Collection
}

func (ctrl *Controller) RegisterHandler(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	// Store user information in MongoDB Atlas
	user.Password = hashedPassword
	_, err = ctrl.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User registration failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func (ctrl *Controller) LoginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := c.ShouldBindJSON(&loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user from MongoDB Atlas by username or email
	var user models.User
	err = ctrl.Collection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"username": loginData.Username}, {"email": loginData.Username}}}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify the hashed password
	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a JWT token
	token, err := ctrl.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ctrl *Controller) GenerateToken(user models.User) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username // Use an appropriate field as the identifier
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

