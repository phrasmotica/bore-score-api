package routes

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

func GetApprovals(c *gin.Context) {
	resultId := c.Param("resultId")

	ctx := context.TODO()

	if !db.ResultExists(ctx, resultId) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("result %s does not exist", resultId)})
		return
	}

	success, approvals := db.GetApprovals(ctx, resultId)

	if !success {
		Error.Println("Could not get approvals")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Got %d approvals\n", len(approvals))

	c.IndentedJSON(http.StatusOK, approvals)
}

func PostApproval(c *gin.Context) {
	var newApproval models.Approval

	if err := c.BindJSON(&newApproval); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if c.GetString("username") != newApproval.Username {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "cannot approve on another user's behalf"})
		c.Abort()
		return
	}

	ctx := context.TODO()

	if !db.UserExists(ctx, newApproval.Username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s does not exist", newApproval.Username)})
		return
	}

	success, result := db.GetResult(ctx, newApproval.ResultID)
	if !success {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("result %s does not exist", newApproval.ResultID)})
		return
	}

	isInResult := false

	for _, score := range result.Scores {
		if newApproval.Username == score.Username {
			isInResult = true
		}
	}

	if !isInResult {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s does not have a score in result %s", newApproval.Username, result.ID)})
		return
	}

	if success := db.AddApproval(ctx, &newApproval); !success {
		Error.Println("Could not add approval")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Added approval for result %s\n", newApproval.ResultID)

	c.IndentedJSON(http.StatusCreated, newApproval)
}
