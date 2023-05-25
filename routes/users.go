package routes

import (
	"context"
	"net/http"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username       string `json:"username" bson:"username"`
	Email          string `json:"email" bson:"email"`
	Password       string `json:"password" bson:"password"`
	DisplayName    string `json:"displayName" bson:"displayName"`
	ProfilePicture string `json:"profilePicture" bson:"profilePicture"`
}

type GetUserResponse struct {
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

	res := &GetUserResponse{
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

	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		Error.Println("Invalid body format")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newUser := models.User{
		ID:          uuid.NewString(),
		Username:    request.Username,
		TimeCreated: time.Now().UTC().Unix(),
		Email:       request.Email,
		Password:    request.Password,
		Permissions: []string{},
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

	success := db.AddUser(ctx, &newUser)
	if !success {
		Error.Printf("Could not add user %s\n", newUser.Username)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Created new user %s\n", newUser.Username)

	// create a player record that corresponds to the new user
	newPlayer := models.Player{
		ID:             uuid.NewString(),
		Username:       newUser.Username,
		TimeCreated:    time.Now().UTC().Unix(),
		DisplayName:    request.DisplayName,
		ProfilePicture: request.ProfilePicture,
	}

	playerSuccess := db.AddPlayer(ctx, &newPlayer)
	if !playerSuccess {
		Error.Printf("Could not add player record for user %s\n", newUser.Username)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Created player record for new user %s\n", newUser.Username)

	c.IndentedJSON(http.StatusNoContent, nil)
}
