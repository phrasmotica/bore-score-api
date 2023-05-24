package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSummary(c *gin.Context) {
	success, summary := db.GetSummary(context.TODO())

	if !success {
		Error.Println("Could not get summary")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Println("Got summary")

	c.IndentedJSON(http.StatusOK, summary)
}
