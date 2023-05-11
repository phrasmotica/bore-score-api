package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/data"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

type ResultResponse struct {
	ID               string                `json:"id" bson:"id"`
	GameName         string                `json:"gameName" bson:"gameName"`
	GroupName        string                `json:"groupName" bson:"groupName"`
	TimeCreated      int64                 `json:"timeCreated" bson:"timeCreated"`
	TimePlayed       int64                 `json:"timePlayed" bson:"timePlayed"`
	Notes            string                `json:"notes" bson:"notes"`
	CooperativeScore int                   `json:"cooperativeScore" bson:"cooperativeScore"`
	CooperativeWin   bool                  `json:"cooperativeWin" bson:"cooperativeWin"`
	Scores           []models.PlayerScore  `json:"scores" bson:"scores"`
	ApprovalStatus   models.ApprovalStatus `json:"approvalStatus" bson:"approvalStatus"`
}

func GetResults(c *gin.Context) {
	username := c.Query("username")
	group := c.Query("group")

	var success bool
	var results []models.Result

	ctx := context.TODO()

	// TODO: allow using both filters simultaneously
	if len(username) > 0 {
		success, results = db.GetResultsWithPlayer(ctx, username)
	} else if len(group) > 0 {
		success, results = db.GetResultsForGroup(ctx, group)
	} else {
		success, results = db.GetAllResults(ctx)
	}

	if !success {
		Error.Println("Could not get results")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	filteredResults := []ResultResponse{}

	callingUsername := c.GetString("username")
	for _, r := range results {
		if canSeeResult(ctx, r, callingUsername) {
			approvalStatus := computeOverallApproval(ctx, db, r)

			filteredResults = append(filteredResults, ResultResponse{
				ID:               r.ID,
				GameName:         r.GameName,
				GroupName:        r.GroupName,
				TimeCreated:      r.TimeCreated,
				TimePlayed:       r.TimePlayed,
				Notes:            r.Notes,
				CooperativeScore: r.CooperativeScore,
				CooperativeWin:   r.CooperativeWin,
				Scores:           r.Scores,
				ApprovalStatus:   approvalStatus,
			})
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

func computeOverallApproval(ctx context.Context, db data.IDatabase, result models.Result) models.ApprovalStatus {
	approvalStatus := models.Pending

	success, approvals := db.GetApprovals(ctx, result.ID)
	if success {
		isApproved := func(a models.Approval) bool { return a.ApprovalStatus == models.Approved }
		isRejected := func(a models.Approval) bool { return a.ApprovalStatus == models.Rejected }

		if all(approvals, isApproved) {
			approvalStatus = models.Approved
		} else if all(approvals, isRejected) {
			approvalStatus = models.Rejected
		}
	} else {
		Error.Println(fmt.Sprintf("Could not get approvals for result %s\n", result.ID))
	}

	return approvalStatus
}

func all[T interface{}](arr []T, predicate func(T) bool) bool {
	for _, e := range arr {
		if !predicate(e) {
			return false
		}
	}

	return true
}
