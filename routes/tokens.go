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
		Error.Println("Invalid body format")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// check if email exists and password is correct
	success, user := db.GetUserByEmail(ctx, request.Email)
	if !success {
		Error.Printf("Could not get user with email %s\n", request.Email)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		Error.Println("Invalid password")
		c.AbortWithError(http.StatusUnauthorized, credentialError)
		return
	}

	tokenString, err := auth.GenerateJWT(user)
	if err != nil {
		Error.Printf("Could not generate token for user with email %s\n", user.Email)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	Info.Printf("Generated token for user %s\n", user.Username)

	c.IndentedJSON(http.StatusOK, gin.H{"token": tokenString})
}

func RefreshToken(c *gin.Context) {
	currentToken := c.GetString("token")
	callingUsername := c.GetString("username")

	tokenString, err := auth.RefreshJWT(currentToken)
	if err != nil {
		Error.Printf("Could not refresh token for user %s\n", callingUsername)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	Info.Printf("Refreshed token for user %s\n", callingUsername)

	c.IndentedJSON(http.StatusOK, gin.H{"token": tokenString})
}
