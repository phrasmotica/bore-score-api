package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetResults(c *gin.Context) {
	username := c.Query("username")

	var success bool
	var results []models.Result

	if len(username) > 0 {
		success, results = db.GetResultsWithPlayer(context.TODO(), username)
	} else {
		success, results = db.GetAllResults(context.TODO())
	}

	if !success {
		Error.Println("Could not get results")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Got %d results\n", len(results))

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

	if len(newResult.GroupName) > 0 && !db.GroupExists(ctx, newResult.GroupName) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("result is attached to non-existent group %s", newResult.GroupName)})
		return
	}

	for _, score := range newResult.Scores {
		if !db.PlayerExists(ctx, score.Username) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s does not exist", score.Username)})
			return
		}
	}

	if success := db.AddResult(ctx, &newResult); !success {
		Error.Println("Could not add result")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Added result for game %s\n", newResult.GameName)

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
