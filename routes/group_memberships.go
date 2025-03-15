package routes

import (
	"context"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

// GetSummary    godoc
// @Summary      Gets a user's group memberships
// @Description  Gets a user's group memberships
// @Tags         Summary
// @Produce      json
// @Param        username path string true "The user's username"
// @Security     BearerAuth
// @Success      200 {object} []models.GroupMembership
// @Failure      401
// @Router       /memberships/{username} [get]
func GetGroupMemberships(c *gin.Context) {
	username := c.Param("username")
	callingUsername := c.GetString("username")

	if username != callingUsername {
		Error.Println("Cannot get another user's group memberships")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, username) {
		Error.Printf("User %s does not exist\n", username)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, memberships := db.GetGroupMemberships(ctx, username)

	if !success {
		Error.Println("Could not get group memberships")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Got %d group memberships\n", len(memberships))

	c.IndentedJSON(http.StatusOK, memberships)
}

func AddGroupMembership(c *gin.Context) {
	var newGroupMembership models.GroupMembership

	if err := c.BindJSON(&newGroupMembership); err != nil {
		Error.Println("Invalid body format")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, newGroupMembership.Username) {
		Error.Printf("User %s does not exist\n", newGroupMembership.Username)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, group := db.GetGroup(ctx, newGroupMembership.GroupID)
	if !success {
		Error.Printf("Group %s does not exist", newGroupMembership.GroupID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if group.Visibility == models.Private {
		Error.Printf("Group %s is private\n", newGroupMembership.GroupID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if db.IsInGroup(ctx, group.ID, newGroupMembership.Username) {
		Info.Printf("User %s is already in group %s\n", newGroupMembership.Username, newGroupMembership.GroupID)
		c.IndentedJSON(http.StatusNoContent, nil)
		return
	}

	if success := db.AddGroupMembership(ctx, &newGroupMembership); !success {
		Error.Println("Could not add group membership")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Added membership to group %s for user %s\n", newGroupMembership.GroupID, newGroupMembership.Username)

	c.IndentedJSON(http.StatusCreated, newGroupMembership)
}
