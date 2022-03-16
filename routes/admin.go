package routes

import (
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/db"

	"github.com/gin-gonic/gin"
)

func GetGameName(c *gin.Context) {
	displayName := c.Param("displayName")
	name := db.ComputeName(displayName)

	fmt.Printf("Computed DB name %s for %s\n", name, displayName)

	c.IndentedJSON(http.StatusOK, gin.H{
		"displayName": displayName,
		"name":        name,
	})
}
