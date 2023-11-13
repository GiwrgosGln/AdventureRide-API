package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"adventureride/models"
	"adventureride/utils"

	"github.com/dgrijalva/jwt-go"
)

// Define secret keys for signing JWT tokens and refresh tokens.
var jwtSecret = []byte("your-secret-key")
var refreshTokenSecret = []byte("your-refresh-token-secret")

type Controller struct {
	Collection *mongo.Collection
}

// CreateIndexes creates the necessary indexes for the MongoDB collection.
func (ctrl *Controller) CreateIndexes() {
	// Create a unique index on the username and email fields
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1, "email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := ctrl.Collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		fmt.Println("Error creating index:", err)
		// Handle the error
	} else {
		fmt.Println("Index created successfully")
	}
}

// TokenDetails holds the access token and refresh token.
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
}

// Constants for token expiration times.
const (
	AccessTokenExpireTime  = time.Hour * 24      // 1 day
	RefreshTokenExpireTime = time.Hour * 24 * 7  // 7 days
)

func (ctrl *Controller) RegisterHandler(c *gin.Context) {
    var user models.User
    err := c.ShouldBindJSON(&user)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if the username or email already exists
    existingUser := models.User{}
    err = ctrl.Collection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"username": user.Username}, {"email": user.Email}}}).Decode(&existingUser)
    if err == nil {
        // User with the same username or email already exists
        c.JSON(http.StatusConflict, gin.H{"error": "Username or email already in use"})
        return
    } else if err != mongo.ErrNoDocuments {
        // Other error occurred
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing user"})
        return
    }

    // Hash the password
    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
        return
    }

    // Store user information in MongoDB
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

    // Generate access and refresh tokens
    tokenDetails, err := ctrl.GenerateTokens(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
        return
    }

    // Include user ID and name in the response
    c.JSON(http.StatusOK, gin.H{
        "access_token":  tokenDetails.AccessToken,
        "refresh_token": tokenDetails.RefreshToken,
        "user_id":       user.ID.Hex(),
        "username":      user.Username,
    })
}


func (ctrl *Controller) GenerateTokens(user models.User) (*TokenDetails, error) {
	// Generate access token
	accessToken, err := ctrl.generateToken(user.Username, user.ID.Hex(), AccessTokenExpireTime)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := ctrl.generateToken(user.Username, user.ID.Hex(), RefreshTokenExpireTime)
	if err != nil {
		return nil, err
	}

	// Store the refresh token in a secure way (e.g., database)

	return &TokenDetails{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ctrl *Controller) generateToken(username, userID string, expirationTime time.Duration) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(expirationTime).Unix()

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}



func (ctrl *Controller) RefreshTokenHandler(c *gin.Context) {
	// Extract the refresh token from the request
	var refreshTokenData struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := c.ShouldBindJSON(&refreshTokenData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Validate the refresh token
	username, err := ctrl.ValidateRefreshToken(refreshTokenData.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate a new access and refresh token pair
	user := models.User{Username: username} // You might need to fetch user details from the database
	tokenDetails, err := ctrl.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tokenDetails.AccessToken, "refresh_token": tokenDetails.RefreshToken})
}

func (ctrl *Controller) ValidateRefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}

		return refreshTokenSecret, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	// Extract username from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("Invalid username in token claims")
	}

	return username, nil
}
