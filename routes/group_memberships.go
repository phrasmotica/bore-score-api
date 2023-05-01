package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/data"
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

	if c.GetString("username") != newGroupMembership.Username {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "cannot join a group on another user's behalf"})
		c.Abort()
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, newGroupMembership.Username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s does not exist", newGroupMembership.Username)})
		return
	}

	result, group := db.GetGroup(ctx, newGroupMembership.GroupID)
	if result != data.Success {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("group %s does not exist", newGroupMembership.GroupID)})
		return
	}

	if group.Visibility == models.Private {
		// TODO: use token to determine who is adding the user to this group
		inviterIsInGroup := db.IsInGroup(ctx, newGroupMembership.GroupID, newGroupMembership.InviterUsername)
		if !inviterIsInGroup {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("inviter is not in group %s", newGroupMembership.GroupID)})
			return
		}
	}

	if success := db.AddGroupMembership(ctx, &newGroupMembership); !success {
		Error.Println("Could not add group membership")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Added membership to group %s for user %s\n", newGroupMembership.GroupID, newGroupMembership.Username)

	c.IndentedJSON(http.StatusCreated, newGroupMembership)
}
