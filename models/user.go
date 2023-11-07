package models

type User struct {
	Username string `bson:"username"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
