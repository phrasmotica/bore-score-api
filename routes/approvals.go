package routes

import (
	"context"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetApprovals(c *gin.Context) {
	resultId := c.Param("resultId")

	ctx := context.TODO()

	if !db.ResultExists(ctx, resultId) {
		Error.Printf("Result %s does not exist\n", resultId)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, approvals := db.GetApprovals(ctx, resultId)

	if !success {
		Error.Println("Could not get approvals")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Got %d approvals\n", len(approvals))

	c.IndentedJSON(http.StatusOK, approvals)
}

func PostApproval(c *gin.Context) {
	var newApproval models.Approval

	if err := c.BindJSON(&newApproval); err != nil {
		Error.Println("Invalid body format")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if c.GetString("username") != newApproval.Username {
		Error.Println("Cannot approve on another user's behalf")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, newApproval.Username) {
		Error.Printf("User %s does not exist", newApproval.Username)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, result := db.GetResult(ctx, newApproval.ResultID)
	if !success {
		Error.Printf("Result %s does not exist", newApproval.ResultID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isInResult := false

	for _, score := range result.Scores {
		// TODO: use slices.ContainsFunc(...) to check
		if newApproval.Username == score.Username {
			isInResult = true
		}
	}

	if !isInResult {
		Error.Printf("Player %s does not have a score in result %s", newApproval.Username, result.ID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if success := db.AddApproval(ctx, &newApproval); !success {
		Error.Println("Could not add approval")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Added approval for result %s\n", newApproval.ResultID)

	c.IndentedJSON(http.StatusCreated, newApproval)
}
