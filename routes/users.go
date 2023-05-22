package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserResponse struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}

func GetUser(c *gin.Context) {
	ctx := context.TODO()

	username := c.Param("username")

	success, user := db.GetUser(ctx, username)

	if !success {
		Error.Printf("Could not get user %s\n", username)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Got user %s\n", username)

	res := &UserResponse{
		Username: username,
	}

	// only return this user's email address if the request was made by this user
	callingUsername := c.GetString("username")
	if callingUsername == username {
		res.Email = user.Email
	}

	c.IndentedJSON(http.StatusOK, res)
}

func RegisterUser(c *gin.Context) {
	ctx := context.TODO()

	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	if err := newUser.HashPassword(newUser.Password); err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	if db.UserExistsByEmail(ctx, newUser.Email) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s already exists", newUser.Email)})
		c.Abort()
		return
	}

	prime(&newUser)

	success := db.AddUser(ctx, &newUser)
	if !success {
		Error.Printf("Could not add user %s\n", newUser.Username)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"email": newUser.Email, "username": newUser.Username})
}

// Primes the user object so that it's ready for database insertion.
func prime(user *models.User) {
	user.ID = uuid.NewString()
	user.TimeCreated = time.Now().UTC().Unix()
}
