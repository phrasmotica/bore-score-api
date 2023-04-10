package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSummary(c *gin.Context) {
	success, summary := db.GetSummary(context.TODO())

	if !success {
		fmt.Println("Could not get summary")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, summary)
}
