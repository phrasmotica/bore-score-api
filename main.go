package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var games = []game{
	{
		ID:         1,
		Name:       "Village Green",
		GameType:   Score,
		MinPlayers: 1,
		MaxPlayers: 5,
	},
	{
		ID:         2,
		Name:       "Modern Art: The Card Game",
		GameType:   Score,
		MinPlayers: 2,
		MaxPlayers: 5,
	},
	{
		ID:         3,
		Name:       "Love Letter",
		GameType:   Score,
		MinPlayers: 2,
		MaxPlayers: 4,
	},
}

var players = []player{
	{
		ID:          1,
		Username:    "johannam",
		DisplayName: "Johanna",
	},
	{
		ID:          2,
		Username:    "julianl",
		DisplayName: "Julian",
	},
}

var results = []result{
	{
		ID:        1,
		GameID:    1,
		Timestamp: time.Date(2022, time.January, 22, 10, 34, 0, 0, time.UTC).Unix(),
		Scores: []playerScore{
			{
				PlayerID: 1,
				Score:    25,
			},
			{
				PlayerID: 2,
				Score:    23,
			},
		},
	},
	{
		ID:        2,
		GameID:    1,
		Timestamp: time.Date(2022, time.January, 23, 17, 12, 0, 0, time.UTC).Unix(),
		Scores: []playerScore{
			{
				PlayerID: 1,
				Score:    32,
			},
			{
				PlayerID: 2,
				Score:    34,
			},
		},
	},
	{
		ID:        3,
		GameID:    2,
		Timestamp: time.Date(2022, time.February, 13, 14, 56, 0, 0, time.UTC).Unix(),
		Scores: []playerScore{
			{
				PlayerID: 1,
				Score:    116,
			},
			{
				PlayerID: 2,
				Score:    140,
			},
		},
	},
}

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/games", getGames)

	router.GET("/players", getPlayers)
	router.POST("/players", postPlayer)
	router.DELETE("/players/:username", deletePlayer)

	router.GET("/results", getResults)
	router.POST("/results", postResult)

	router.Run("localhost:8000")
}

func getGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, games)
}

func getPlayers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, players)
}

func postPlayer(c *gin.Context) {
	var newPlayer player

	if err := c.BindJSON(&newPlayer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if playerExistsByUsername(players, newPlayer.Username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s already exists", newPlayer.Username)})
		return
	}

	newPlayer.ID = getMaxPlayerId(players) + 1

	players = append(players, newPlayer)
	c.IndentedJSON(http.StatusCreated, newPlayer)
}

func deletePlayer(c *gin.Context) {
	username := c.Param("username")

	if !playerExistsByUsername(players, username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s does not exist", username)})
		return
	}

	players = removePlayer(players, username)
	c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func getResults(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, results)
}

func postResult(c *gin.Context) {
	var newResult result

	if err := c.BindJSON(&newResult); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if !gameExists(games, newResult.GameID) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %d does not exist", newResult.GameID)})
		return
	}

	for _, score := range newResult.Scores {
		if !playerExists(players, score.PlayerID) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %d does not exist", score.PlayerID)})
			return
		}
	}

	newResult.ID = getMaxResultId(results) + 1

	results = append(results, newResult)
	c.IndentedJSON(http.StatusCreated, newResult)
}
