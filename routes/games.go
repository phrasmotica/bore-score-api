package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"phrasmotica/bore-score-api/data"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

// TODO: put these in a more central place, or inject them as dependencies
var (
	Info  *log.Logger = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile)
	Error *log.Logger = log.New(os.Stdout, "ERROR: ", log.LstdFlags|log.Lshortfile)
)

// TODO: put this in a more central place, or inject it as a dependency
func createDb() data.IDatabase {
	azureTablesConnStr := os.Getenv("AZURE_TABLES_CONNECTION_STRING")
	if azureTablesConnStr != "" {
		Info.Println("Using data backend: Azure Table Storage")

		return &data.TableStorageDatabase{
			Client: data.CreateTableStorageClient(azureTablesConnStr),
		}
	}

	mongoDbUri := os.Getenv("MONGODB_URI")
	if mongoDbUri != "" {
		Info.Println("Using data backend: MongoDB")

		return &data.MongoDatabase{
			Database: data.CreateMongoDatabase(mongoDbUri),
		}
	}

	panic("No AZURE_TABLES_CONNECTION_STRING or MONGODB_URI environment variable found!")
}

var db = createDb()

func GetGames(c *gin.Context) {
	success, games := db.GetAllGames(context.TODO())

	if !success {
		Error.Println("Could not get games")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Got %d games\n", len(games))

	c.IndentedJSON(http.StatusOK, games)
}

func GetGame(c *gin.Context) {
	name := c.Param("name")

	success, game := db.GetGame(context.TODO(), name)

	if !success {
		Error.Printf("Game %s does not exist\n", name)
		c.IndentedJSON(http.StatusNotFound, gin.H{})
		return
	}

	Info.Printf("Got game %s\n", name)

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
		Error.Printf("Could not add game %s\n", newGame.Name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Added game %s\n", newGame.Name)

	c.IndentedJSON(http.StatusCreated, newGame)
}

func validateNewGame(game *models.Game) (bool, string) {
	if len(game.Name) <= 0 {
		return false, "game name is missing"
	}

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

	success, deletedCount := db.DeleteResultsWithGame(ctx, name)
	if !success {
		Error.Printf("Could not delete results for game %s\n", name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Deleted %d results for game %s\n", deletedCount, name)

	if success := db.DeleteGame(ctx, name); !success {
		Error.Printf("Could not delete game %s\n", name)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Deleted game %s\n", name)

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}
