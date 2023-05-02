package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetAllGroups(c *gin.Context) {
	ctx := context.TODO()

	success, groups := db.GetAllGroups(ctx)

	if !success {
		Error.Println("Could not get all groups")
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

func GetGroups(c *gin.Context) {
	ctx := context.TODO()

	success, groups := db.GetGroups(ctx)

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
	name := c.Param("name")

	ctx := context.TODO()

	success, group := db.GetGroupByName(ctx, name)

	if !success {
		Error.Printf("Group %s not found\n", name)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "group not found"})
		return
	}

	if group.Visibility == models.Private {
		callingUsername := c.GetString("username")
		isInGroup := db.IsInGroup(ctx, group.ID, callingUsername)

		if !isInGroup {
			Error.Printf("Group %s is private\n", name)
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "group is private"})
			return
		}
	}

	Info.Printf("Got group %s\n", name)
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

	if db.GroupExists(ctx, newGroup.Name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("group %s already exists", newGroup.Name)})
		return
	}

	if success := db.AddGroup(ctx, &newGroup); !success {
		Error.Printf("Could not add group %s\n", newGroup.Name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Added group %s\n", newGroup.Name)

	c.IndentedJSON(http.StatusCreated, newGroup)
}

func validateNewGroup(group *models.Group) (bool, string) {
	if len(group.DisplayName) <= 0 {
		return false, "group display name is missing"
	}

	return true, ""
}

func DeleteGroup(c *gin.Context) {
	name := c.Param("name")

	ctx := context.TODO()

	if !db.GroupExists(ctx, name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("group %s does not exist", name)})
		return
	}

	if success := db.DeleteGroup(ctx, name); !success {
		Error.Printf("Could not delete group %s\n", name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Deleted group %s\n", name)

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func canSeeGroup(ctx context.Context, group models.Group, callingUsername string) bool {
	return group.Visibility != models.Private || db.IsInGroup(ctx, group.ID, callingUsername)
}
