package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	ID       primitive.ObjectID `bson:"_id"`
}
