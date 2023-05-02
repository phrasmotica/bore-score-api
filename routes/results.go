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

	ctx := context.TODO()

	if len(username) > 0 {
		success, results = db.GetResultsWithPlayer(ctx, username)
	} else {
		success, results = db.GetAllResults(ctx)
	}

	if !success {
		Error.Println("Could not get results")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	filteredResults := []models.Result{}

	callingUsername := c.GetString("username")
	for _, r := range results {
		if canSeeResult(ctx, r, callingUsername) {
			filteredResults = append(filteredResults, r)
		}
	}

	Info.Printf("Got %d results\n", len(filteredResults))

	c.IndentedJSON(http.StatusOK, filteredResults)
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

	if len(newResult.GroupName) > 0 {
		success, group := db.GetGroupByName(ctx, newResult.GroupName)
		if !success {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("result is attached to non-existent group %s", newResult.GroupName)})
			return
		}

		for _, score := range newResult.Scores {
			if !db.IsInGroup(ctx, group.ID, score.Username) {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s is not in group %s", score.Username, newResult.GroupName)})
				return
			}
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

func canSeeResult(ctx context.Context, r models.Result, callingUsername string) bool {
	if len(r.GroupName) <= 0 {
		return true
	}

	success, group := db.GetGroupByName(ctx, r.GroupName)
	return success && canSeeGroup(ctx, *group, callingUsername)
}
