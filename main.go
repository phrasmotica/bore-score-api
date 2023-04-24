package main

import (
	"phrasmotica/bore-score-api/auth"
	"phrasmotica/bore-score-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	games := router.Group("/games")
	{
		games.GET("", routes.GetGames)
		games.GET("/:name", routes.GetGame)
		games.POST("", routes.PostGame)
		games.DELETE("/:name", routes.DeleteGame)
	}

	groups := router.Group("/groups")
	{
		groups.GET("", routes.GetGroups)

		// TODO: use a route param instead of a separate route
		groups.GET("-all", routes.GetAllGroups)

		groups.GET("/:name", routes.GetGroup)
		groups.POST("", routes.PostGroup)
		groups.DELETE("/:name", routes.DeleteGroup)
	}

	linkTypes := router.Group("/linkTypes")
	{
		linkTypes.GET("", routes.GetLinkTypes)
	}

	players := router.Group("/players")
	{
		players.GET("", routes.GetPlayers)
		players.GET("/:username", routes.GetPlayer)
		players.POST("", routes.PostPlayer)
		players.DELETE("/:username", routes.DeletePlayer)
	}

	router.GET("/summary", routes.GetSummary)

	results := router.Group("/results")
	{
		results.GET("", routes.GetResults)
		results.POST("", routes.PostResult)
	}

	winMethods := router.Group("/winMethods")
	{
		winMethods.GET("", routes.GetWinMethods)
	}

	token := router.Group("/token")
	{
		token.POST("", routes.GenerateToken)
	}

	users := router.Group("/users")
	{
		users.GET("/:username", auth.TokenAuth(true), routes.GetUser)
		users.POST("", routes.RegisterUser)
	}

	secured := router.Group("/secured").Use(auth.TokenAuth(false))
	{
		secured.GET("/ping", routes.Ping)
	}

	router.Run(":8000")
}
