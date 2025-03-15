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

type UpdatePasswordRequest struct {
	Username        string `json:"username" bson:"username"`
	CurrentPassword string `json:"currentPassword" bson:"currentPassword"`
	NewPassword     string `json:"newPassword" bson:"newPassword"`
}

// GetSummary    godoc
// @Summary      Gets a user
// @Description  Gets a user
// @Tags         Summary
// @Produce      json
// @Param        username path string true "The user's username"
// @Security     BearerAuth
// @Success      200 {object} routes.GetUserResponse
// @Router       /users/{username} [get]
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

func UpdatePassword(c *gin.Context) {
	username := c.Param("username")
	callingUsername := c.GetString("username")

	var request UpdatePasswordRequest

	ctx := context.TODO()

	if err := c.BindJSON(&request); err != nil {
		Error.Println("Invalid body format")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if success, err := validateUpdatePasswordRequest(&request); !success {
		Error.Printf("Error validating update password request: %s\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if username != request.Username {
		Error.Println("Update password request is for wrong user")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if callingUsername != request.Username {
		Error.Println("Cannot update password of a different user")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	exists, user := db.GetUser(ctx, request.Username)
	if !exists {
		Error.Printf("User %s does not exist", request.Username)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	credentialError := user.CheckPassword(request.CurrentPassword)
	if credentialError != nil {
		Error.Println("Invalid password")
		c.AbortWithError(http.StatusUnauthorized, credentialError)
		return
	}

	if err := user.HashPassword(request.NewPassword); err != nil {
		Error.Println("Could not hash password")
		c.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}

	success := db.UpdateUser(ctx, user)
	if !success {
		Error.Printf("Could not update password for user %s\n", request.Username)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Updated password for user %s\n", request.Username)

	c.IndentedJSON(http.StatusNoContent, nil)
}

func validateUpdatePasswordRequest(request *UpdatePasswordRequest) (bool, string) {
	if len(request.Username) <= 0 {
		return false, "username is missing"
	}

	if len(request.CurrentPassword) <= 0 {
		return false, "current password is missing"
	}

	if len(request.NewPassword) <= 0 {
		return false, "new password is missing"
	}

	return true, ""
}
