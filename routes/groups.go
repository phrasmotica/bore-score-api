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

	filteredGroups := []models.Group{}

	for _, g := range groups {
		callingUsername := c.GetString("username")

		if canSeeGroup(ctx, g, callingUsername) {
			filteredGroups = append(filteredGroups, g)
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

		if !db.IsInGroup(ctx, group.ID, callingUsername) {
			Error.Printf("User %s is not in group %s\n", callingUsername, groupId)
			c.IndentedJSON(http.StatusUnauthorized, gin.H{})
			return
		}
	}

	Info.Printf("Got group %s\n", groupId)

	c.IndentedJSON(http.StatusOK, group)
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

func canSeeGroup(ctx context.Context, group models.Group, callingUsername string) bool {
	return group.Visibility != models.Private || db.IsInGroup(ctx, group.ID, callingUsername)
}
