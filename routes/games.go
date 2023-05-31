package routes

import (
	"context"
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
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Got %d games\n", len(games))

	c.IndentedJSON(http.StatusOK, games)
}

func GetGame(c *gin.Context) {
	id := c.Param("gameId")

	success, game := db.GetGame(context.TODO(), id)

	if !success {
		Error.Printf("Game %s does not exist\n", id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	Info.Printf("Got game %s\n", id)

	c.IndentedJSON(http.StatusOK, game)
}

func PostGame(c *gin.Context) {
	var newGame models.Game

	ctx := context.TODO()

	if err := c.BindJSON(&newGame); err != nil {
		Error.Println("Invalid body format")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if success, err := validateNewGame(&newGame); !success {
		Error.Printf("Error validating new game %s: %s\n", newGame.DisplayName, err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if success := db.AddGame(ctx, &newGame); !success {
		Error.Printf("Could not add game %s\n", newGame.ID)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Added game %s\n", newGame.ID)

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
	id := c.Param("gameId")

	ctx := context.TODO()

	if !db.GameExists(ctx, id) {
		Error.Printf("Game %s does not exist", id)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, deletedCount := db.DeleteResultsWithGame(ctx, id)
	if !success {
		Error.Printf("Could not delete results for game %s\n", id)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Deleted %d results for game %s\n", deletedCount, id)

	if success := db.DeleteGame(ctx, id); !success {
		Error.Printf("Could not delete game %s\n", id)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Deleted game %s\n", id)

	c.IndentedJSON(http.StatusNoContent, nil)
}
