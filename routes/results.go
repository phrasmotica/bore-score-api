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

func GetResults(c *gin.Context) {
	results, success := db.GetAllResults(context.TODO())

	if !success {
		fmt.Println("Could not get results")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got %d results\n", len(results))

	c.IndentedJSON(http.StatusOK, results)
}

func PostResult(c *gin.Context) {
	var newResult models.Result

	if err := c.BindJSON(&newResult); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if success, err := validateNewResult(&newResult); !success {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
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

	success := db.AddResult(ctx, &newResult)
	if !success {
		log.Println("Could not add result")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Added result for game %s\n", newResult.GameName)

	c.IndentedJSON(http.StatusCreated, newResult)
}

func validateNewResult(result *models.Result) (bool, string) {
	if len(result.Scores) <= 0 {
		return false, "result is missing player scores"
	}

	if !hasUniquePlayerScores(result) {
		return false, "result has duplicated player scores"
	}

	return true, ""
}
