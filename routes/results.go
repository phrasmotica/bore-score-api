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
	var success bool
	var results []models.Result

	ctx := context.TODO()

	success, results = db.GetAllResults(ctx)
	if !success {
		Error.Println("Could not get results")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	callingUsername := c.GetString("username")
	filteredResults := filterResults(ctx, results, callingUsername)

	Info.Printf("Got %d results\n", len(filteredResults))

	c.IndentedJSON(http.StatusOK, filteredResults)
}

func GetResultsForGroup(c *gin.Context) {
	groupId := c.Param("groupId")

	ctx := context.TODO()

	groupSuccess, group := db.GetGroup(ctx, groupId)
	if !groupSuccess {
		Error.Printf("Group %s does not exist\n", groupId)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	callingUsername := c.GetString("username")

	if !canSeeGroup(ctx, group, callingUsername, false) {
		Error.Printf("User %s cannot see results for group %s\n", callingUsername, group.ID)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	success, results := db.GetResultsForGroup(ctx, groupId)
	if !success {
		Error.Printf("Could not get results for group %s\n", groupId)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	resultResponses := []ResultResponse{}

	for _, r := range results {
		resultResponses = append(resultResponses, createResultResponse(ctx, &r))
	}

	Info.Printf("Got %d results\n", len(resultResponses))

	c.IndentedJSON(http.StatusOK, resultResponses)
}

func GetResultsForUser(c *gin.Context) {
	username := c.Param("username")
	callingUsername := c.GetString("username")

	if username != callingUsername {
		Error.Println("Cannot see results for another user")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := context.TODO()

	success, results := db.GetResultsWithPlayer(ctx, username)
	if !success {
		Error.Printf("Could not get results for user %s\n", username)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	filteredResults := filterResults(ctx, results, callingUsername)

	Info.Printf("Got %d results\n", len(filteredResults))

	c.IndentedJSON(http.StatusOK, filteredResults)
}

func PostResult(c *gin.Context) {
	var newResult models.Result

	if err := c.BindJSON(&newResult); err != nil {
		Error.Println("Invalid body format")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if success, err := validateNewResult(&newResult); !success {
		Error.Printf("Error validating new result %s: %s\n", newResult.ID, err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := context.TODO()

	if !db.GameExists(ctx, newResult.GameName) {
		Error.Printf("Game %s does not exist\n", newResult.GameName)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	for _, score := range newResult.Scores {
		if !db.PlayerExists(ctx, score.Username) {
			Error.Printf("Player %s does not exist\n", score.Username)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	if len(newResult.GroupID) > 0 {
		success, group := db.GetGroup(ctx, newResult.GroupID)
		if !success {
			Error.Printf("Group %s does not exist\n", newResult.GroupID)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		for _, score := range newResult.Scores {
			if !db.IsInGroup(ctx, group.ID, score.Username) {
				Error.Printf("Player %s is not in group %s\n", score.Username, newResult.GroupID)
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
	}

	if success := db.AddResult(ctx, &newResult); !success {
		Error.Println("Could not add result")
		c.AbortWithStatus(http.StatusServiceUnavailable)
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
			filteredResults = append(filteredResults, createResultResponse(ctx, &r))
		}
	}

	return filteredResults
}

func createResultResponse(ctx context.Context, result *models.Result) ResultResponse {
	approvalStatus := computeOverallApproval(ctx, db, result)

	return ResultResponse{
		ID:               result.ID,
		GameName:         result.GameName,
		GroupID:          result.GroupID,
		TimeCreated:      result.TimeCreated,
		TimePlayed:       result.TimePlayed,
		Notes:            result.Notes,
		CooperativeScore: result.CooperativeScore,
		CooperativeWin:   result.CooperativeWin,
		Scores:           result.Scores,
		ApprovalStatus:   approvalStatus,
	}
}

func canSeeResult(ctx context.Context, r models.Result, callingUsername string) bool {
	if len(r.GroupID) <= 0 {
		return true
	}

	success, group := db.GetGroup(ctx, r.GroupID)
	return success && canSeeGroup(ctx, group, callingUsername, false)
}

func computeOverallApproval(ctx context.Context, db data.IDatabase, result *models.Result) models.ApprovalStatus {
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
