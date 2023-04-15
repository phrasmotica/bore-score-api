package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLinkTypes(c *gin.Context) {
	success, linkTypes := db.GetAllLinkTypes(context.TODO())

	if !success {
		fmt.Println("Could not get link types")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got %d link types\n", len(linkTypes))

	c.IndentedJSON(http.StatusOK, linkTypes)
}
