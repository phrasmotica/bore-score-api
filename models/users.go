package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string `json:"id" bson:"id"`
	Username    string `json:"username" bson:"username"`
	TimeCreated int64  `json:"timeCreated" bson:"timeCreated"`
	Email       string `json:"email" bson:"email"`
	Password    string `json:"password" bson:"password"`
}
