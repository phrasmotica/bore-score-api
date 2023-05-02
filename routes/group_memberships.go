package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetGroupMemberships(c *gin.Context) {
	username := c.Param("username")
	callingUsername := c.GetString("username")

	if username != callingUsername {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "cannot get another user's group memberships"})
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s does not exist", username)})
		return
	}

	success, approvals := db.GetGroupMemberships(ctx, username)

	if !success {
		Error.Println("Could not get group memberships")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Got %d group memberships\n", len(approvals))

	c.IndentedJSON(http.StatusOK, approvals)
}

func AddGroupMembership(c *gin.Context) {
	var newGroupMembership models.GroupMembership

	if err := c.BindJSON(&newGroupMembership); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, newGroupMembership.Username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s does not exist", newGroupMembership.Username)})
		return
	}

	success, group := db.GetGroup(ctx, newGroupMembership.GroupID)
	if !success {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("group %s does not exist", newGroupMembership.GroupID)})
		return
	}

	if group.Visibility == models.Private {
		inviterUsername := c.GetString("username")

		inviterIsInGroup := db.IsInGroup(ctx, newGroupMembership.GroupID, inviterUsername)
		if !inviterIsInGroup {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("inviter is not in group %s", newGroupMembership.GroupID)})
			return
		}

		newGroupMembership.InviterUsername = inviterUsername
	}

	if db.IsInGroup(ctx, group.ID, newGroupMembership.Username) {
		Info.Printf("User %s is already in group %s\n", newGroupMembership.Username, newGroupMembership.GroupID)
		c.IndentedJSON(http.StatusNoContent, gin.H{})
		return
	}

	if success := db.AddGroupMembership(ctx, &newGroupMembership); !success {
		Error.Println("Could not add group membership")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Added membership to group %s for user %s by inviter %s\n", newGroupMembership.GroupID, newGroupMembership.Username, newGroupMembership.InviterUsername)

	c.IndentedJSON(http.StatusCreated, newGroupMembership)
}
