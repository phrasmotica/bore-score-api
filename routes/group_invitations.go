package routes

import (
	"context"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetGroupInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	success, invitation := db.GetGroupInvitation(context.TODO(), invitationId)

	if !success {
		Error.Printf("Group invitation %s does not exist\n", invitationId)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	callingUsername := c.GetString("username")

	if invitation.InviterUsername != callingUsername {
		Error.Println("Cannot get another user's group invitations")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	Info.Printf("Got group invitation %s\n", invitationId)

	c.IndentedJSON(http.StatusOK, invitation)
}

func GetGroupInvitationsForUser(c *gin.Context) {
	username := c.Param("username")
	callingUsername := c.GetString("username")

	if username != callingUsername {
		Error.Println("Cannot get another user's group invitations")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, username) {
		Info.Printf("User %s does not exist", username)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, approvals := db.GetGroupInvitations(ctx, username)

	if !success {
		Error.Println("Could not get group invitations")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Got %d group invitations\n", len(approvals))

	c.IndentedJSON(http.StatusOK, approvals)
}

func GetGroupInvitationsForGroup(c *gin.Context) {
	groupId := c.Param("groupId")

	ctx := context.TODO()

	success, group := db.GetGroup(ctx, groupId)
	if !success {
		Error.Printf("Group %s does not exist\n", groupId)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	callingUsername := c.GetString("username")

	if !db.IsInGroup(ctx, group.ID, callingUsername) {
		Error.Printf("User %s is not in group %s\n", callingUsername, groupId)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	success, invitations := db.GetGroupInvitationsForGroup(ctx, groupId)

	if !success {
		Error.Printf("Could not get group invitations for group %s\n", groupId)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Got %d group invitations\n", len(invitations))

	c.IndentedJSON(http.StatusOK, invitations)
}

func AddGroupInvitation(c *gin.Context) {
	var newGroupInvitation models.GroupInvitation

	if err := c.BindJSON(&newGroupInvitation); err != nil {
		Error.Println("Invalid body format")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, newGroupInvitation.Username) {
		Error.Printf("User %s does not exist\n", newGroupInvitation.Username)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !db.UserExists(ctx, newGroupInvitation.InviterUsername) {
		Error.Printf("User %s does not exist\n", newGroupInvitation.InviterUsername)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, group := db.GetGroup(ctx, newGroupInvitation.GroupID)
	if !success {
		Error.Printf("Group %s does not exist\n", newGroupInvitation.GroupID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !db.IsInGroup(ctx, newGroupInvitation.GroupID, newGroupInvitation.InviterUsername) {
		Error.Printf("Inviter %s is not in group %s\n", newGroupInvitation.InviterUsername, newGroupInvitation.GroupID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if db.IsInGroup(ctx, group.ID, newGroupInvitation.Username) {
		Info.Printf("User %s is already in group %s\n", newGroupInvitation.Username, newGroupInvitation.GroupID)
		c.IndentedJSON(http.StatusNoContent, nil)
		return
	}

	if success := db.AddGroupInvitation(ctx, &newGroupInvitation); !success {
		Error.Println("Could not add group invitation")
		c.AbortWithStatus(http.StatusServiceUnavailable)
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
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	callingUsername := c.GetString("username")

	if invitation.Username != callingUsername {
		Error.Println("Cannot accept another user's group invitation")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if !db.UserExists(ctx, invitation.Username) {
		Error.Printf("Invited user %s does not exist\n", invitation.Username)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if !db.UserExists(ctx, invitation.InviterUsername) {
		Error.Printf("Inviting user %s does not exist\n", invitation.InviterUsername)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if !db.IsInGroup(ctx, invitation.GroupID, invitation.InviterUsername) {
		Error.Printf("Inviter %s is not in group %s\n", invitation.InviterUsername, invitation.GroupID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if db.IsInGroup(ctx, invitation.GroupID, invitation.Username) {
		Info.Printf("User %s is already in group %s\n", invitation.Username, invitation.GroupID)
		c.IndentedJSON(http.StatusNoContent, nil)
		return
	}

	invitation.InvitationStatus = models.Accepted

	if success := db.UpdateGroupInvitation(ctx, invitation); !success {
		Error.Println("Could not accept group invitation")
		c.AbortWithStatus(http.StatusServiceUnavailable)
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
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Added membership to group %s for user %s from invitation %s\n", newMembership.GroupID, newMembership.Username, newMembership.InvitationID)

	c.IndentedJSON(http.StatusNoContent, nil)
}

func DeclineGroupInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	ctx := context.TODO()

	success, invitation := db.GetGroupInvitation(ctx, invitationId)
	if !success {
		Error.Printf("Group invitation %s does not exist\n", invitationId)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	callingUsername := c.GetString("username")

	if invitation.Username != callingUsername {
		Error.Println("Cannot decline another user's group invitation")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if !db.UserExists(ctx, invitation.Username) {
		Error.Printf("Invited user %s does not exist\n", invitation.Username)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if !db.UserExists(ctx, invitation.InviterUsername) {
		Error.Printf("Inviting user %s does not exist\n", invitation.InviterUsername)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if !db.IsInGroup(ctx, invitation.GroupID, invitation.InviterUsername) {
		Error.Printf("Inviter %s is not in group %s\n", invitation.InviterUsername, invitation.GroupID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if db.IsInGroup(ctx, invitation.GroupID, invitation.Username) {
		Info.Printf("User %s is already in group %s\n", invitation.Username, invitation.GroupID)
		c.IndentedJSON(http.StatusNoContent, nil)
		return
	}

	invitation.InvitationStatus = models.Declined

	if success := db.UpdateGroupInvitation(ctx, invitation); !success {
		Error.Println("Could not decline group invitation")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Declined invitation to group %s for user %s by inviter %s\n", invitation.GroupID, invitation.Username, invitation.InviterUsername)

	c.IndentedJSON(http.StatusNoContent, nil)
}
