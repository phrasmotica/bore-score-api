package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLinkTypes(c *gin.Context) {
	success, linkTypes := db.GetAllLinkTypes(context.TODO())

	if !success {
		Error.Println("Could not get link types")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Got %d link types\n", len(linkTypes))

	c.IndentedJSON(http.StatusOK, linkTypes)
}
