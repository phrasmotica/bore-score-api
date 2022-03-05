package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/db"

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
