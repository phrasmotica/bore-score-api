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

func GetGames(c *gin.Context) {
	games, success := db.GetAllGames(context.TODO())

	if !success {
		fmt.Println("Could not get games")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got %d games\n", len(games))

	c.IndentedJSON(http.StatusOK, games)
}

func GetGame(c *gin.Context) {
	name := c.Param("name")

	game, success := db.GetGame(context.TODO(), name)

	if !success {
		fmt.Printf("Could not get game %s\n", name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got game %s\n", name)

	c.IndentedJSON(http.StatusOK, game)
}

func PostGame(c *gin.Context) {
	var newGame models.Game

	ctx := context.TODO()

	if err := c.BindJSON(&newGame); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if success, err := validateNewGame(&newGame); !success {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	if db.GameExists(ctx, newGame.Name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %s already exists", newGame.Name)})
		return
	}

	if success := db.AddGame(ctx, &newGame); !success {
		log.Printf("Could not add game %s\n", newGame.Name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Added game %s\n", newGame.Name)

	c.IndentedJSON(http.StatusCreated, newGame)
}

func validateNewGame(game *models.Game) (bool, string) {
	if len(game.DisplayName) <= 0 {
		return false, "game display name is missing"
	}

	if game.MinPlayers <= 0 {
		return false, "game min players must be at least 1"
	}

	if game.MaxPlayers < game.MinPlayers {
		return false, "game max players must be at least equal to its min players"
	}

	if len(game.WinMethod) <= 0 {
		return false, "game display name is missing"
	}

	return true, ""
}

func DeleteGame(c *gin.Context) {
	name := c.Param("name")

	ctx := context.TODO()

	if !db.GameExists(ctx, name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %s does not exist", name)})
		return
	}

	deletedCount, success := db.DeleteResultsWithGame(ctx, name)
	if !success {
		log.Printf("Could not delete results for game %s\n", name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Deleted %d results for game %s\n", deletedCount, name)

	if success := db.DeleteGame(ctx, name); !success {
		log.Printf("Could not delete game %s\n", name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Deleted game %s\n", name)

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}
