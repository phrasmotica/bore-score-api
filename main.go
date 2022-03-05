package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"phrasmotica/bore-score-api/db"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/summary", getSummary)

	router.GET("/games", getGames)
	router.GET("/games/:name", getGame)
	router.POST("/games", postGame)
	router.DELETE("/games/:name", deleteGame)

	router.GET("/linkTypes", getLinkTypes)

	router.GET("/winMethods", getWinMethods)

	router.GET("/groups", getGroups)
	router.GET("/groups/:name", getGroup)

	router.GET("/players", getPlayers)
	router.POST("/players", postPlayer)
	router.DELETE("/players/:username", deletePlayer)

	router.GET("/results", getResults)
	router.POST("/results", postResult)

	router.Run("localhost:8000")
}

func getSummary(c *gin.Context) {
	summary := db.GetSummary(context.TODO())

	c.IndentedJSON(http.StatusOK, summary)
}

func getGames(c *gin.Context) {
	games := db.GetAllGames(context.TODO())

	fmt.Printf("Found %d games\n", len(games))

	c.IndentedJSON(http.StatusOK, games)
}

func getGame(c *gin.Context) {
	name := c.Param("name")

	game := db.GetGame(context.TODO(), name)

	fmt.Printf("Found game %s\n", name)

	c.IndentedJSON(http.StatusOK, game)
}

func postGame(c *gin.Context) {
	var newGame models.Game

	ctx := context.TODO()

	if err := c.BindJSON(&newGame); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if err := validateNewGame(&newGame); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if db.GameExists(ctx, newGame.Name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %s already exists", newGame.Name)})
		return
	}

	err := db.AddGame(ctx, &newGame)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added game %s\n", newGame.Name)

	c.IndentedJSON(http.StatusCreated, newGame)
}

func validateNewGame(game *models.Game) error {
	if len(game.DisplayName) <= 0 {
		return errors.New("game display name is missing")
	}

	if game.MinPlayers <= 0 {
		return errors.New("game min players must be at least 1")
	}

	if game.MaxPlayers < game.MinPlayers {
		return errors.New("game max players must be at least equal to its min players")
	}

	if len(game.WinMethod) <= 0 {
		return errors.New("game display name is missing")
	}

	return nil
}

func deleteGame(c *gin.Context) {
	name := c.Param("name")

	ctx := context.TODO()

	if !db.GameExists(ctx, name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %s does not exist", name)})
		return
	}

	deletedCount, err := db.DeleteResultsWithGame(ctx, name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted %d results for game %s\n", deletedCount, name)

	_, err = db.DeleteGame(ctx, name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted game %s\n", name)

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func getWinMethods(c *gin.Context) {
	winMethods := db.GetAllWinMethods(context.TODO())

	fmt.Printf("Found %d win methods\n", len(winMethods))

	c.IndentedJSON(http.StatusOK, winMethods)
}

func getLinkTypes(c *gin.Context) {
	linkTypes := db.GetAllLinkTypes(context.TODO())

	fmt.Printf("Found %d link types\n", len(linkTypes))

	c.IndentedJSON(http.StatusOK, linkTypes)
}

func getGroups(c *gin.Context) {
	groups := db.GetAllGroups(context.TODO())

	fmt.Printf("Found %d groups\n", len(*groups))

	c.IndentedJSON(http.StatusOK, groups)
}

func getGroup(c *gin.Context) {
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

	fmt.Printf("Found group %s\n", name)
	c.IndentedJSON(http.StatusOK, group)
}

func getPlayers(c *gin.Context) {
	players := db.GetAllPlayers(context.TODO())

	fmt.Printf("Found %d players\n", len(players))

	c.IndentedJSON(http.StatusOK, players)
}

func postPlayer(c *gin.Context) {
	var newPlayer models.Player

	ctx := context.TODO()

	if err := c.BindJSON(&newPlayer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if err := validateNewPlayer(&newPlayer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if db.PlayerExists(ctx, newPlayer.Username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s already exists", newPlayer.Username)})
		return
	}

	err := db.AddPlayer(ctx, &newPlayer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added player %s\n", newPlayer.Username)

	c.IndentedJSON(http.StatusCreated, newPlayer)
}

func validateNewPlayer(player *models.Player) error {
	if len(player.Username) <= 0 {
		return errors.New("player username is missing")
	}

	if len(player.DisplayName) <= 0 {
		return errors.New("player display name is missing")
	}

	return nil
}

func deletePlayer(c *gin.Context) {
	username := c.Param("username")

	ctx := context.TODO()

	if !db.PlayerExists(ctx, username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s does not exist", username)})
		return
	}

	scrubbedCount, err := db.ScrubResultsWithPlayer(ctx, username)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Scrubbed player %s from %d results\n", username, scrubbedCount)

	_, err = db.DeletePlayer(ctx, username)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted player %s\n", username)

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func getResults(c *gin.Context) {
	results := db.GetAllResults(context.TODO())

	fmt.Printf("Found %d results\n", len(results))

	c.IndentedJSON(http.StatusOK, results)
}

func postResult(c *gin.Context) {
	var newResult models.Result

	if err := c.BindJSON(&newResult); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if err := validateNewResult(&newResult); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx := context.TODO()

	if !db.GameExists(ctx, newResult.GameName) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %s does not exist", newResult.GameName)})
		return
	}

	for _, score := range newResult.Scores {
		if !db.PlayerExists(ctx, score.Username) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s does not exist", score.Username)})
			return
		}
	}

	err := db.AddResult(ctx, &newResult)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added result for game %s\n", newResult.GameName)

	c.IndentedJSON(http.StatusCreated, newResult)
}

func validateNewResult(result *models.Result) error {
	if len(result.Scores) <= 0 {
		return errors.New("result is missing player scores")
	}

	if !hasUniquePlayerScores(result) {
		return errors.New("result has duplicated player scores")
	}

	return nil
}
