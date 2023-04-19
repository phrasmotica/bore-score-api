package main

import (
	"phrasmotica/bore-score-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/games", routes.GetGames)
	router.GET("/games/:name", routes.GetGame)
	router.POST("/games", routes.PostGame)
	router.DELETE("/games/:name", routes.DeleteGame)

	router.GET("/groups", routes.GetGroups)
	router.GET("/groups-all", routes.GetAllGroups)
	router.GET("/groups/:name", routes.GetGroup)
	router.POST("/groups", routes.PostGroup)
	router.DELETE("/groups/:name", routes.DeleteGroup)

	router.GET("/linkTypes", routes.GetLinkTypes)

	router.GET("/players", routes.GetPlayers)
	router.GET("/players/:username", routes.GetPlayer)
	router.POST("/players", routes.PostPlayer)
	router.DELETE("/players/:username", routes.DeletePlayer)

	router.GET("/summary", routes.GetSummary)

	router.GET("/results", routes.GetResults)
	router.POST("/results", routes.PostResult)

	router.GET("/winMethods", routes.GetWinMethods)

	router.POST("/user/register", routes.RegisterUser)
	router.POST("/token", routes.GenerateToken)

	router.Run(":8000")
}
