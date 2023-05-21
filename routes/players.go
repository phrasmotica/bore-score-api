package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetPlayers(c *gin.Context) {
	success, players := db.GetAllPlayers(context.TODO())
	if !success {
		Error.Println("Could not get players")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{})
		return
	}

	Info.Printf("Got %d players\n", len(players))

	c.IndentedJSON(http.StatusOK, players)
}

func GetPlayersInGroup(c *gin.Context) {
	groupId := c.Param("groupId")

	ctx := context.TODO()

	groupSuccess, group := db.GetGroup(ctx, groupId)
	if !groupSuccess {
		Error.Printf("Group %s does not exist\n", groupId)
		c.IndentedJSON(http.StatusNotFound, gin.H{})
		return
	}

	callingUsername := c.GetString("username")

	if !canSeeGroup(ctx, group, callingUsername, false) {
		Error.Printf("User %s cannot see results for group %s\n", callingUsername, group.ID)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	success, players := db.GetPlayersInGroup(ctx, group.ID)
	if !success {
		Error.Printf("Could not get players in group %s\n", group.ID)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{})
		return
	}

	Info.Printf("Got %d players\n", len(players))

	c.IndentedJSON(http.StatusOK, players)
}

func GetPlayer(c *gin.Context) {
	username := c.Param("username")

	success, player := db.GetPlayer(context.TODO(), username)

	if !success {
		Error.Printf("Player %s does not exist\n", username)
		c.IndentedJSON(http.StatusNotFound, gin.H{})
		return
	}

	Info.Printf("Got player %s\n", username)

	c.IndentedJSON(http.StatusOK, player)
}

func PostPlayer(c *gin.Context) {
	var newPlayer models.Player

	ctx := context.TODO()

	if err := c.BindJSON(&newPlayer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if success, err := validateNewPlayer(&newPlayer); !success {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	if db.PlayerExists(ctx, newPlayer.Username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s already exists", newPlayer.Username)})
		return
	}

	success := db.AddPlayer(ctx, &newPlayer)
	if !success {
		Error.Printf("Could not add player %s\n", newPlayer.Username)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Added player %s\n", newPlayer.Username)

	c.IndentedJSON(http.StatusCreated, newPlayer)
}

func validateNewPlayer(player *models.Player) (bool, string) {
	if len(player.Username) <= 0 {
		return false, "player username is missing"
	}

	if len(player.DisplayName) <= 0 {
		return false, "player display name is missing"
	}

	return true, ""
}

func DeletePlayer(c *gin.Context) {
	username := c.Param("username")

	ctx := context.TODO()

	if !db.PlayerExists(ctx, username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s does not exist", username)})
		return
	}

	success, scrubbedCount := db.ScrubResultsWithPlayer(ctx, username)
	if !success {
		Error.Printf("Could not scrub player %s from results\n", username)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Scrubbed player %s from %d results\n", username, scrubbedCount)

	if success := db.DeletePlayer(ctx, username); !success {
		Error.Printf("Could not delete player %s\n", username)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Deleted player %s\n", username)

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}
