package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWinMethods(c *gin.Context) {
	success, winMethods := db.GetAllWinMethods(context.TODO())

	if !success {
		Error.Println("Could not get win methods")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	Info.Printf("Got %d win methods\n", len(winMethods))

	c.IndentedJSON(http.StatusOK, winMethods)
}
