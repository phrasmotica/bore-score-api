package routes

import (
	"context"
	"net/http"
	"phrasmotica/bore-score-api/auth"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateToken(c *gin.Context) {
	ctx := context.TODO()

	var request TokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	// check if email exists and password is correct
	success, user := db.GetUserByEmail(ctx, request.Email)
	if !success {
		Error.Printf("Could not get user with email %s\n", request.Email)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		c.Abort()
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		c.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	Info.Printf("Generated token for user %s\n", user.Username)

	c.IndentedJSON(http.StatusOK, gin.H{"token": tokenString})
}

func RefreshToken(c *gin.Context) {
	currentToken := c.GetString("token")

	tokenString, err := auth.RefreshJWT(currentToken)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	Info.Printf("Refreshed token for user %s\n", c.GetString("username"))

	c.IndentedJSON(http.StatusOK, gin.H{"token": tokenString})
}
