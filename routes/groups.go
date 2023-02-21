package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"phrasmotica/bore-score-api/db"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetAllGroups(c *gin.Context) {
	groups, success := db.GetAllGroups(context.TODO())

	if !success {
		fmt.Println("Could not get all groups")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got %d groups\n", len(groups))

	c.IndentedJSON(http.StatusOK, groups)
}

func GetGroups(c *gin.Context) {
	groups, success := db.GetGroups(context.TODO())

	if !success {
		fmt.Println("Could not get groups")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got %d groups\n", len(groups))

	c.IndentedJSON(http.StatusOK, groups)
}

func GetGroup(c *gin.Context) {
	name := c.Param("name")

	group, result := db.GetGroup(context.TODO(), name)

	if result == db.Failure {
		fmt.Printf("Group %s not found\n", name)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "group not found"})
		return
	}

	if result == db.Unauthorised {
		fmt.Printf("Group %s is private\n", name)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "group is private"})
		return
	}

	fmt.Printf("Got group %s\n", name)
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
		log.Printf("Could not add group %s\n", newGroup.Name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Added group %s\n", newGroup.Name)

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
		log.Printf("Could not delete group %s\n", name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Deleted group %s\n", name)

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}