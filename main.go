package main

import (
	"context"
	"fmt"
	"net/http"
	"phrasmotica/bore-score-api/db"
	"phrasmotica/bore-score-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/summary", getSummary)
	router.GET("/linkTypes", getLinkTypes)
	router.GET("/winMethods", getWinMethods)

	router.GET("/games", routes.GetGames)
	router.GET("/games/:name", routes.GetGame)
	router.POST("/games", routes.PostGame)
	router.DELETE("/games/:name", routes.DeleteGame)

	router.GET("/groups", routes.GetGroups)
	router.GET("/groups-all", routes.GetAllGroups)
	router.GET("/groups/:name", routes.GetGroup)

	router.GET("/players", routes.GetPlayers)
	router.GET("/players/:username", routes.GetPlayer)
	router.POST("/players", routes.PostPlayer)
	router.DELETE("/players/:username", routes.DeletePlayer)

	router.GET("/results", routes.GetResults)
	router.POST("/results", routes.PostResult)

	router.GET("/admin/game-name/:displayName", routes.GetGameName)

	router.Run("localhost:8000")
}

func getSummary(c *gin.Context) {
	summary, success := db.GetSummary(context.TODO())

	if !success {
		fmt.Println("Could not get summary")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, summary)
}

func getLinkTypes(c *gin.Context) {
	linkTypes, success := db.GetAllLinkTypes(context.TODO())

	if !success {
		fmt.Println("Could not get link types")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got %d link types\n", len(linkTypes))

	c.IndentedJSON(http.StatusOK, linkTypes)
}

func getWinMethods(c *gin.Context) {
	winMethods, success := db.GetAllWinMethods(context.TODO())

	if !success {
		fmt.Println("Could not get win methods")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got %d win methods\n", len(winMethods))

	c.IndentedJSON(http.StatusOK, winMethods)
}
