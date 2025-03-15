package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSummary    godoc
// @Summary      Gets a summary of the database
// @Description  Gets a summary of the database
// @Tags         Summary
// @Produce      json
// @Success      200 {object} data.Summary
// @Router       /summary [get]
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
