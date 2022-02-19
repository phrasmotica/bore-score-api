package main

import (
	"context"
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

	router.GET("/games", getGames)
	router.POST("/games", postGame)
	router.DELETE("/games/:id", deleteGame)

	router.GET("/winMethods", getWinMethods)

	router.GET("/players", getPlayers)
	router.POST("/players", postPlayer)
	router.DELETE("/players/:username", deletePlayer)

	router.GET("/results", getResults)
	router.POST("/results", postResult)

	router.Run("localhost:8000")
}

func getGames(c *gin.Context) {
	games := db.GetAllGames(context.TODO())

	fmt.Printf("Found %d games\n", len(games))

	c.IndentedJSON(http.StatusOK, games)
}

func postGame(c *gin.Context) {
	var newGame models.Game

	ctx := context.TODO()

	if err := c.BindJSON(&newGame); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if db.GameExists(ctx, newGame.Name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %d already exists", newGame.ID)})
		return
	}

	err := db.AddGame(ctx, &newGame)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added game %s\n", newGame.Name)

	c.IndentedJSON(http.StatusCreated, newGame)
}

func deleteGame(c *gin.Context) {
	name := c.Param("name")

	ctx := context.TODO()

	if !db.GameExists(ctx, name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %s does not exist", name)})
		return
	}

	deletedCount, err := db.DeleteResultsWithGameId(ctx, name)
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
	fmt.Printf("Found %d win methods\n", len(models.WinMethods))

	c.IndentedJSON(http.StatusOK, models.WinMethods)
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

	ctx := context.TODO()

	if err := c.BindJSON(&newResult); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if !db.GameExists(ctx, newResult.GameName) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %d does not exist", newResult.GameName)})
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
