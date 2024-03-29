package routes

import (
	"context"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetPlayers(c *gin.Context) {
	success, players := db.GetAllPlayers(context.TODO())
	if !success {
		Error.Println("Could not get players")
		c.AbortWithStatus(http.StatusServiceUnavailable)
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
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	callingUsername := c.GetString("username")

	if !canSeeGroup(ctx, group, callingUsername, false) {
		Error.Printf("User %s cannot see results for group %s\n", callingUsername, group.ID)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	success, players := db.GetPlayersInGroup(ctx, group.ID)
	if !success {
		Error.Printf("Could not get players in group %s\n", group.ID)
		c.AbortWithStatus(http.StatusServiceUnavailable)
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
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	Info.Printf("Got player %s\n", username)

	c.IndentedJSON(http.StatusOK, player)
}

func PostPlayer(c *gin.Context) {
	var newPlayer models.Player

	ctx := context.TODO()

	if err := c.BindJSON(&newPlayer); err != nil {
		Error.Println("Invalid body format")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if success, err := validatePlayer(&newPlayer); !success {
		Error.Printf("Error validating new player %s: %s\n", newPlayer.DisplayName, err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if db.PlayerExists(ctx, newPlayer.Username) {
		Error.Printf("Player %s already exists", newPlayer.Username)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success := db.AddPlayer(ctx, &newPlayer)
	if !success {
		Error.Printf("Could not add player %s\n", newPlayer.Username)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Added player %s\n", newPlayer.Username)

	c.IndentedJSON(http.StatusCreated, newPlayer)
}

func validatePlayer(player *models.Player) (bool, string) {
	if len(player.Username) <= 0 {
		return false, "player username is missing"
	}

	if len(player.DisplayName) <= 0 {
		return false, "player display name is missing"
	}

	return true, ""
}

func UpdatePlayer(c *gin.Context) {
	username := c.Param("username")
	callingUsername := c.GetString("username")

	var player models.Player

	ctx := context.TODO()

	if err := c.BindJSON(&player); err != nil {
		Error.Println("Invalid body format")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO: only need profile picture and display name
	if success, err := validatePlayer(&player); !success {
		Error.Printf("Error validating player %s: %s\n", player.DisplayName, err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if username != player.Username {
		Error.Println("Update request is for wrong player")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if callingUsername != player.Username {
		Error.Println("Cannot update a different player")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	exists, existingPlayer := db.GetPlayer(ctx, player.Username)
	if !exists {
		Error.Printf("Player %s does not exist", player.Username)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// ensure we update the correct entity
	player.ID = existingPlayer.ID

	success := db.UpdatePlayer(ctx, &player)
	if !success {
		Error.Printf("Could not update player %s\n", player.Username)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Updated player %s\n", player.Username)

	c.IndentedJSON(http.StatusNoContent, nil)
}

func DeletePlayer(c *gin.Context) {
	username := c.Param("username")

	ctx := context.TODO()

	if !db.PlayerExists(ctx, username) {
		Error.Printf("Player %s does not exist\n", username)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, scrubbedCount := db.ScrubResultsWithPlayer(ctx, username)
	if !success {
		Error.Printf("Could not scrub player %s from results\n", username)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Scrubbed player %s from %d results\n", username, scrubbedCount)

	if success := db.DeletePlayer(ctx, username); !success {
		Error.Printf("Could not delete player %s\n", username)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Deleted player %s\n", username)

	c.IndentedJSON(http.StatusNoContent, nil)
}
