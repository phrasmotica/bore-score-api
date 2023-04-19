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

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}
