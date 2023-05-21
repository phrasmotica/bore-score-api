package routes

import (
	"context"
	"net/http"
	"phrasmotica/bore-score-api/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupResponse struct {
	ID             string                     `json:"id" bson:"id"`
	TimeCreated    int64                      `json:"timeCreated" bson:"timeCreated"`
	DisplayName    string                     `json:"displayName" bson:"displayName"`
	Description    string                     `json:"description" bson:"description"`
	ProfilePicture string                     `json:"profilePicture" bson:"profilePicture"`
	CreatedBy      string                     `json:"createdBy" bson:"createdBy"`
	Visibility     models.GroupVisibilityName `json:"visibility" bson:"visibility"`
	MemberCount    int                        `json:"memberCount" bson:"memberCount"`
}

func GetGroups(c *gin.Context) {
	getAll := c.Query("all") == strconv.Itoa(1)

	var success bool
	var groups []models.Group

	ctx := context.TODO()

	if getAll {
		success, groups = db.GetAllGroups(ctx)
	} else {
		success, groups = db.GetGroups(ctx)
	}

	if !success {
		Error.Println("Could not get groups")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	filteredGroups := []GroupResponse{}

	for _, g := range groups {
		callingUsername := c.GetString("username")

		if canSeeGroup(ctx, &g, callingUsername, true) {
			filteredGroups = append(filteredGroups, createGroupResponse(ctx, &g))
		}
	}

	Info.Printf("Got %d groups\n", len(filteredGroups))

	c.IndentedJSON(http.StatusOK, filteredGroups)
}

func GetGroup(c *gin.Context) {
	groupId := c.Param("groupId")

	ctx := context.TODO()

	success, group := db.GetGroup(ctx, groupId)

	if !success {
		Error.Printf("Group %s not found\n", groupId)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "group not found"})
		return
	}

	if group.Visibility == models.Private {
		callingUsername := c.GetString("username")

		// could use canSeeGroup(...) here, but prefer to break down the conditions
		// for logging purposes. TODO: put logging into canSeeGroup(...)?
		if !db.IsInGroup(ctx, group.ID, callingUsername) {
			if !db.IsInvitedToGroup(ctx, group.ID, callingUsername) {
				Error.Printf("User %s is not in private group %s\n", callingUsername, groupId)
				c.IndentedJSON(http.StatusUnauthorized, gin.H{})
				return
			} else {
				Info.Printf("User %s is invited to private group %s\n", callingUsername, group.ID)
			}
		}
	}

	groupResponse := createGroupResponse(ctx, group)

	Info.Printf("Got group %s\n", groupResponse.ID)

	c.IndentedJSON(http.StatusOK, groupResponse)
}

func PostGroup(c *gin.Context) {
	var newGroup models.Group

	ctx := context.TODO()

	if err := c.BindJSON(&newGroup); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if success, err := validateNewGroup(&newGroup); !success {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	creatorUsername := c.GetString("username")
	newGroup.CreatedBy = creatorUsername

	newGroup.ID = uuid.NewString()
	newGroup.TimeCreated = time.Now().UTC().Unix()

	if success := db.AddGroup(ctx, &newGroup); !success {
		Error.Printf("Could not add group %s\n", newGroup.DisplayName)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Added group %s\n", newGroup.DisplayName)

	// add membership for the creator
	membership := models.GroupMembership{
		ID:           uuid.NewString(),
		GroupID:      newGroup.ID,
		TimeCreated:  time.Now().UTC().Unix(),
		Username:     creatorUsername,
		InvitationID: "",
	}

	if success := db.AddGroupMembership(ctx, &membership); !success {
		// not a fatal error, they can join the group afterwards...
		Error.Printf("Could not add membership to group %s for group creator %s\n", newGroup.ID, creatorUsername)
	} else {
		Info.Printf("Added membership to group %s for group creator %s\n", newGroup.ID, creatorUsername)
	}

	c.IndentedJSON(http.StatusCreated, newGroup)
}

func validateNewGroup(group *models.Group) (bool, string) {
	if len(group.DisplayName) <= 0 {
		return false, "group display name is missing"
	}

	return true, ""
}

func DeleteGroup(c *gin.Context) {
	groupId := c.Param("groupId")

	ctx := context.TODO()

	success, group := db.GetGroup(ctx, groupId)
	if !success {
		Error.Printf("Group %s does not exist\n", groupId)
		c.IndentedJSON(http.StatusNotFound, gin.H{})
		return
	}

	callingUsername := c.GetString("username")
	if group.CreatedBy != callingUsername {
		Error.Println("Cannot delete a group that someone else created")
		c.IndentedJSON(http.StatusForbidden, gin.H{})
		return
	}

	if success := db.DeleteGroup(ctx, group.ID); !success {
		Error.Printf("Could not delete group %s\n", groupId)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{})
		return
	}

	Info.Printf("Deleted group %s\n", groupId)

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func createGroupResponse(ctx context.Context, group *models.Group) GroupResponse {
	memberCount := computeMemberCount(ctx, group)

	return GroupResponse{
		ID:             group.ID,
		TimeCreated:    group.TimeCreated,
		DisplayName:    group.DisplayName,
		Description:    group.Description,
		ProfilePicture: group.ProfilePicture,
		CreatedBy:      group.CreatedBy,
		Visibility:     group.Visibility,
		MemberCount:    memberCount,
	}
}

func canSeeGroup(ctx context.Context, group *models.Group, callingUsername string, allowInvitees bool) bool {
	return (group.Visibility != models.Private ||
		db.IsInGroup(ctx, group.ID, callingUsername) ||
		(allowInvitees && db.IsInvitedToGroup(ctx, group.ID, callingUsername)))
}

func computeMemberCount(ctx context.Context, group *models.Group) int {
	success, members := db.GetPlayersInGroup(ctx, group.ID)
	if !success {
		return 0
	}

	return len(members)
}
