package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetGroupInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	success, invitation := db.GetGroupInvitation(context.TODO(), invitationId)

	if !success {
		Error.Printf("Group invitation %s does not exist\n", invitationId)
		c.IndentedJSON(http.StatusNotFound, nil)
		return
	}

	callingUsername := c.GetString("username")

	if invitation.InviterUsername != callingUsername {
		Error.Println("Cannot get another user's group invitations")
		c.IndentedJSON(http.StatusUnauthorized, nil)
		return
	}

	Info.Printf("Got group invitation %s\n", invitationId)

	c.IndentedJSON(http.StatusOK, invitation)
}

func GetGroupInvitations(c *gin.Context) {
	username := c.Param("username")
	callingUsername := c.GetString("username")

	if username != callingUsername {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "cannot get another user's group invitations"})
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s does not exist", username)})
		return
	}

	success, approvals := db.GetGroupInvitations(ctx, username)

	if !success {
		Error.Println("Could not get group invitations")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Got %d group invitations\n", len(approvals))

	c.IndentedJSON(http.StatusOK, approvals)
}

func AddGroupInvitation(c *gin.Context) {
	var newGroupInvitation models.GroupInvitation

	if err := c.BindJSON(&newGroupInvitation); err != nil {
		Error.Println("Invalid body format")
		c.IndentedJSON(http.StatusBadRequest, gin.H{})
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, newGroupInvitation.Username) {
		Error.Printf("User %s does not exist\n", newGroupInvitation.Username)
		c.IndentedJSON(http.StatusBadRequest, gin.H{})
		return
	}

	if !db.UserExists(ctx, newGroupInvitation.InviterUsername) {
		Error.Printf("User %s does not exist\n", newGroupInvitation.InviterUsername)
		c.IndentedJSON(http.StatusBadRequest, gin.H{})
		return
	}

	success, group := db.GetGroup(ctx, newGroupInvitation.GroupID)
	if !success {
		Error.Printf("Group %s does not exist\n", newGroupInvitation.GroupID)
		c.IndentedJSON(http.StatusBadRequest, gin.H{})
		return
	}

	if !db.IsInGroup(ctx, newGroupInvitation.GroupID, newGroupInvitation.InviterUsername) {
		Error.Printf("Inviter %s is not in group %s\n", newGroupInvitation.InviterUsername, newGroupInvitation.GroupID)
		c.IndentedJSON(http.StatusForbidden, gin.H{})
		return
	}

	if db.IsInGroup(ctx, group.ID, newGroupInvitation.Username) {
		Info.Printf("User %s is already in group %s\n", newGroupInvitation.Username, newGroupInvitation.GroupID)
		c.IndentedJSON(http.StatusNoContent, gin.H{})
		return
	}

	if success := db.AddGroupInvitation(ctx, &newGroupInvitation); !success {
		Error.Println("Could not add group invitation")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{})
		return
	}

	Info.Printf("Added invitation to group %s for user %s by inviter %s\n", newGroupInvitation.GroupID, newGroupInvitation.Username, newGroupInvitation.InviterUsername)

	c.IndentedJSON(http.StatusCreated, newGroupInvitation)
}

func AcceptGroupInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	ctx := context.TODO()

	success, invitation := db.GetGroupInvitation(ctx, invitationId)
	if !success {
		Error.Printf("Group invitation %s does not exist\n", invitationId)
		c.IndentedJSON(http.StatusForbidden, gin.H{})
		return
	}

	callingUsername := c.GetString("username")

	if invitation.Username != callingUsername {
		Error.Println("Cannot accept another user's group invitation")
		c.IndentedJSON(http.StatusForbidden, gin.H{})
		return
	}

	if !db.UserExists(ctx, invitation.Username) {
		Error.Printf("Invited user %s does not exist\n", invitation.Username)
		c.IndentedJSON(http.StatusForbidden, gin.H{})
		return
	}

	if !db.UserExists(ctx, invitation.InviterUsername) {
		Error.Printf("Inviting user %s does not exist\n", invitation.InviterUsername)
		c.IndentedJSON(http.StatusForbidden, gin.H{})
		return
	}

	if !db.IsInGroup(ctx, invitation.GroupID, invitation.InviterUsername) {
		Error.Printf("Inviter %s is not in group %s\n", invitation.InviterUsername, invitation.GroupID)
		c.IndentedJSON(http.StatusForbidden, gin.H{})
		return
	}

	if db.IsInGroup(ctx, invitation.GroupID, invitation.Username) {
		Info.Printf("User %s is already in group %s\n", invitation.Username, invitation.GroupID)
		c.IndentedJSON(http.StatusNoContent, gin.H{})
		return
	}

	invitation.InvitationStatus = models.Accepted

	if success := db.UpdateGroupInvitation(ctx, invitation); !success {
		Error.Println("Could not accept group invitation")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{})
		return
	}

	Info.Printf("Accepted invitation to group %s for user %s by inviter %s\n", invitation.GroupID, invitation.Username, invitation.InviterUsername)

	newMembership := models.GroupMembership{
		GroupID:      invitation.GroupID,
		Username:     invitation.Username,
		InvitationID: invitation.ID,
	}

	if success := db.AddGroupMembership(ctx, &newMembership); !success {
		Error.Println("Could not add group membership ")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{})
		return
	}

	Info.Printf("Added membership to group %s for user %s from invitation %s\n", newMembership.GroupID, newMembership.Username, newMembership.InvitationID)

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}
