package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

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

	if db.UserExists(ctx, newUser.Email) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s already exists", newUser.Email)})
		c.Abort()
		return
	}

	success := db.AddUser(ctx, &newUser)
	if !success {
		Error.Printf("Could not add user %s\n", newUser.Username)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"email": newUser.Email, "username": newUser.Username})
}
