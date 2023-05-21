package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/data"
	"phrasmotica/bore-score-api/models"
	"sort"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

type ResultResponse struct {
	ID               string                `json:"id" bson:"id"`
	GameName         string                `json:"gameName" bson:"gameName"`
	GroupID          string                `json:"groupId" bson:"groupId"`
	TimeCreated      int64                 `json:"timeCreated" bson:"timeCreated"`
	TimePlayed       int64                 `json:"timePlayed" bson:"timePlayed"`
	Notes            string                `json:"notes" bson:"notes"`
	CooperativeScore int                   `json:"cooperativeScore" bson:"cooperativeScore"`
	CooperativeWin   bool                  `json:"cooperativeWin" bson:"cooperativeWin"`
	Scores           []models.PlayerScore  `json:"scores" bson:"scores"`
	ApprovalStatus   models.ApprovalStatus `json:"approvalStatus" bson:"approvalStatus"`
}

func GetResults(c *gin.Context) {
	groupId := c.Query("group")

	callingUsername := c.GetString("username")

	var success bool
	var results []models.Result

	ctx := context.TODO()

	if len(groupId) > 0 {
		// TODO: move this to a new endpoint /groups/{groupId}/results
		groupSuccess, group := db.GetGroup(ctx, groupId)

		if !groupSuccess {
			Error.Printf("Group %s does not exist\n", groupId)
			c.IndentedJSON(http.StatusNotFound, gin.H{})
			return
		}

		if !canSeeGroup(ctx, group, callingUsername, false) {
			Error.Printf("User %s cannot see results for group %s\n", callingUsername, group.ID)
			c.IndentedJSON(http.StatusUnauthorized, gin.H{})
			return
		}

		success, results = db.GetResultsForGroup(ctx, groupId)
	} else {
		success, results = db.GetAllResults(ctx)
	}

	if !success {
		Error.Println("Could not get results")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	filteredResults := filterResults(ctx, results, callingUsername)

	Info.Printf("Got %d results\n", len(filteredResults))

	c.IndentedJSON(http.StatusOK, filteredResults)
}

func GetResultsForUser(c *gin.Context) {
	username := c.Param("username")
	callingUsername := c.GetString("username")

	if username != callingUsername {
		Error.Println("Cannot see results for another user")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx := context.TODO()

	success, results := db.GetResultsWithPlayer(ctx, username)
	if !success {
		Error.Printf("Could not get results for user %s\n", username)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{})
		return
	}

	filteredResults := filterResults(ctx, results, callingUsername)

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

	if len(newResult.GroupID) > 0 {
		success, group := db.GetGroup(ctx, newResult.GroupID)
		if !success {
			Error.Printf("Result is attached to non-existent group %s", newResult.GroupID)
			c.IndentedJSON(http.StatusForbidden, gin.H{})
			return
		}

		for _, score := range newResult.Scores {
			if !db.IsInGroup(ctx, group.ID, score.Username) {
				Error.Printf("Player %s is not in group %s", score.Username, newResult.GroupID)
				c.IndentedJSON(http.StatusForbidden, gin.H{})
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

func filterResults(ctx context.Context, results []models.Result, username string) []ResultResponse {
	filteredResults := []ResultResponse{}

	for _, r := range results {
		if canSeeResult(ctx, r, username) {
			approvalStatus := computeOverallApproval(ctx, db, r)

			filteredResults = append(filteredResults, ResultResponse{
				ID:               r.ID,
				GameName:         r.GameName,
				GroupID:          r.GroupID,
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

	return filteredResults
}

func canSeeResult(ctx context.Context, r models.Result, callingUsername string) bool {
	if len(r.GroupID) <= 0 {
		return true
	}

	success, group := db.GetGroup(ctx, r.GroupID)
	return success && canSeeGroup(ctx, group, callingUsername, false)
}

func computeOverallApproval(ctx context.Context, db data.IDatabase, result models.Result) models.ApprovalStatus {
	approvalStatus := models.Pending

	success, approvals := db.GetApprovals(ctx, result.ID)
	if success {
		isApproved := func(a models.Approval) bool { return a.ApprovalStatus == models.Approved }
		isRejected := func(a models.Approval) bool { return a.ApprovalStatus == models.Rejected }

		latestApprovals := computeLatestApprovals(approvals)

		if len(latestApprovals) == len(result.Scores) {
			if all(latestApprovals, isApproved) {
				approvalStatus = models.Approved
			} else if all(latestApprovals, isRejected) {
				approvalStatus = models.Rejected
			}
		}
	} else {
		Error.Println(fmt.Sprintf("Could not get approvals for result %s\n", result.ID))
	}

	return approvalStatus
}

func computeLatestApprovals(approvals []models.Approval) []models.Approval {
	latestApprovals := []models.Approval{}

	sortedApprovals := approvals[:]

	// https://stackoverflow.com/a/42872183
	sort.Slice(sortedApprovals, func(i, j int) bool {
		return sortedApprovals[i].TimeCreated > sortedApprovals[j].TimeCreated
	})

	usersAdded := []string{}

	for _, a := range sortedApprovals {
		if !slices.Contains(usersAdded, a.Username) {
			latestApprovals = append(latestApprovals, a)
			usersAdded = append(usersAdded, a.Username)
		}
	}

	return latestApprovals
}

func all[T interface{}](arr []T, predicate func(T) bool) bool {
	for _, e := range arr {
		if !predicate(e) {
			return false
		}
	}

	return true
}
