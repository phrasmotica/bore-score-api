package routes

import (
	"context"
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
		c.AbortWithStatus(http.StatusServiceUnavailable)
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
		Error.Println("Invalid body format")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := newUser.HashPassword(newUser.Password); err != nil {
		Error.Println("Could not hash password")
		c.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}

	if db.UserExistsByEmail(ctx, newUser.Email) {
		Error.Printf("User %s already exists\n", newUser.Email)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	prime(&newUser)

	success := db.AddUser(ctx, &newUser)
	if !success {
		Error.Printf("Could not add user %s\n", newUser.Username)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"email": newUser.Email, "username": newUser.Username})
}

// Primes the user object so that it's ready for database insertion.
func prime(user *models.User) {
	user.ID = uuid.NewString()
	user.TimeCreated = time.Now().UTC().Unix()
	user.Permissions = []string{}
}
