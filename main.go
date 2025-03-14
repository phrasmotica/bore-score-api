package main

import (
	"phrasmotica/bore-score-api/auth"
	docs "phrasmotica/bore-score-api/docs/borescoreapi"
	"phrasmotica/bore-score-api/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// adapted from https://levelup.gitconnected.com/tutorial-generate-swagger-specification-and-swaggerui-for-gin-go-web-framework-9f0c038483b5, https://github.com/swaggo/gin-swagger

// @title BoreScore API
// @version 0.1.0
// @description This is the BoreScore API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
// @schemes http
func main() {
	router := gin.Default()

	router.Use(auth.CORSMiddleware())

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	approvals := router.Group("/approvals", auth.TokenAuth(false))
	{
		approvals.GET("/:resultId", routes.GetApprovals)

		approvals.POST("", routes.PostApproval)
	}

	games := router.Group("/games")
	{
		games.GET("", routes.GetGames)

		games.POST("", routes.PostGame)

		gameByName := games.Group("/:gameId")
		{
			gameByName.GET("", routes.GetGame)

			gameByName.DELETE("", auth.TokenAuth(false), auth.CheckPermission("superuser"), routes.DeleteGame)
		}
	}

	groups := router.Group("/groups")
	{
		groups.GET("", auth.TokenAuth(true), routes.GetGroups)

		groups.POST("", auth.TokenAuth(false), routes.PostGroup)

		groupById := groups.Group("/:groupId")
		{
			groupById.GET("", auth.TokenAuth(true), routes.GetGroup)
			groupById.GET("/invitations", auth.TokenAuth(false), routes.GetGroupInvitationsForGroup)
			groupById.GET("/players", auth.TokenAuth(false), routes.GetPlayersInGroup)
			groupById.GET("/results", auth.TokenAuth(false), routes.GetResultsForGroup)

			groupById.DELETE("", auth.TokenAuth(false), auth.CheckPermission("superuser"), routes.DeleteGroup)

			groupLeaderboards := groupById.Group("/leaderboard")
			{
				groupLeaderboards.GET("/:gameId", auth.TokenAuth(false), routes.GetLeaderboard)
			}
		}
	}

	groupInvitations := router.Group("/invitations", auth.TokenAuth(false))
	{
		groupInvitations.POST("", routes.AddGroupInvitation)

		groupInvitationById := groupInvitations.Group("/:invitationId")
		{
			groupInvitationById.GET("", routes.GetGroupInvitation)

			groupInvitationById.POST("/accept", routes.AcceptGroupInvitation)
			groupInvitationById.POST("/decline", routes.DeclineGroupInvitation)
		}
	}

	groupMemberships := router.Group("/memberships", auth.TokenAuth(false))
	{
		groupMemberships.GET("/:username", routes.GetGroupMemberships)

		groupMemberships.POST("", routes.AddGroupMembership)
	}

	linkTypes := router.Group("/linkTypes")
	{
		linkTypes.GET("", routes.GetLinkTypes)
	}

	// TODO: move Player columns into User entity
	players := router.Group("/players")
	{
		players.GET("", routes.GetPlayers)

		players.POST("", routes.PostPlayer)

		playerByUsername := players.Group("/:username")
		{
			playerByUsername.GET("", routes.GetPlayer)

			playerByUsername.PUT("", auth.TokenAuth(false), routes.UpdatePlayer)

			playerByUsername.DELETE("", auth.TokenAuth(false), auth.CheckPermission("superuser"), routes.DeletePlayer)
		}
	}

	router.GET("/summary", routes.GetSummary)

	results := router.Group("/results")
	{
		results.GET("", auth.TokenAuth(true), routes.GetResults)

		results.POST("", routes.PostResult)
	}

	winMethods := router.Group("/winMethods")
	{
		winMethods.GET("", routes.GetWinMethods)
	}

	token := router.Group("/token")
	{
		token.POST("", routes.GenerateToken)
		token.POST("/refresh", auth.TokenAuth(false), routes.RefreshToken)
	}

	users := router.Group("/users")
	{
		users.POST("", routes.RegisterUser)

		userByUsername := users.Group("/:username")
		{
			userByUsername.GET("", auth.TokenAuth(true), routes.GetUser)
			userByUsername.GET("/invitations", auth.TokenAuth(false), routes.GetGroupInvitationsForUser)
			userByUsername.GET("/results", auth.TokenAuth(false), routes.GetResultsForUser)

			userByUsername.PUT("/password", auth.TokenAuth(false), routes.UpdatePassword)
		}
	}

	router.Run(":8000")
}
