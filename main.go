package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"phrasmotica/bore-score-api/data"
	"phrasmotica/bore-score-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var db data.IDatabase

func main() {
	loadEnv()

	db = createDatabase()

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
	router.POST("/groups", routes.PostGroup)
	router.DELETE("/groups/:name", routes.DeleteGroup)

	router.GET("/players", routes.GetPlayers)
	router.GET("/players/:username", routes.GetPlayer)
	router.POST("/players", routes.PostPlayer)
	router.DELETE("/players/:username", routes.DeletePlayer)

	router.GET("/results", routes.GetResults)
	router.POST("/results", routes.PostResult)

	router.Run(":8000")
}

func loadEnv() {
	env := os.Getenv("BORESCORE_ENV")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	godotenv.Load()
}

func createDatabase() data.IDatabase {
	return &data.MongoDatabase{
		Database: data.CreateMongoDatabase(),
	}
}

func getSummary(c *gin.Context) {
	success, summary := db.GetSummary(context.TODO())

	if !success {
		fmt.Println("Could not get summary")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, summary)
}

func getLinkTypes(c *gin.Context) {
	success, linkTypes := db.GetAllLinkTypes(context.TODO())

	if !success {
		fmt.Println("Could not get link types")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got %d link types\n", len(linkTypes))

	c.IndentedJSON(http.StatusOK, linkTypes)
}

func getWinMethods(c *gin.Context) {
	success, winMethods := db.GetAllWinMethods(context.TODO())

	if !success {
		fmt.Println("Could not get win methods")
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "something went wrong"})
		return
	}

	fmt.Printf("Got %d win methods\n", len(winMethods))

	c.IndentedJSON(http.StatusOK, winMethods)
}
