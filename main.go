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
	router.GET("/swagger/index.html", ginSwagger.WrapHandler(swaggerFiles.Handler))

	approvals := router.Group("/approvals").Use(auth.TokenAuth(false))
	{
		approvals.GET("/:resultId", routes.GetApprovals)
		approvals.POST("", routes.PostApproval)
	}

	games := router.Group("/games")
	{
		games.GET("", routes.GetGames)
		games.GET("/:name", routes.GetGame)
		games.POST("", routes.PostGame)
		games.DELETE("/:name", routes.DeleteGame)
	}

	groups := router.Group("/groups")
	{
		groups.GET("", auth.TokenAuth(true), routes.GetGroups)

		groups.GET("/:groupId", auth.TokenAuth(true), routes.GetGroup)
		groups.GET("/:groupId/invitations", auth.TokenAuth(false), routes.GetGroupInvitationsForGroup)

		groups.POST("", auth.TokenAuth(false), routes.PostGroup)

		groups.DELETE("/:groupId", auth.TokenAuth(false), routes.DeleteGroup)
	}

	groupInvitations := router.Group("/invitations").Use(auth.TokenAuth(false))
	{
		groupInvitations.GET("/:invitationId", routes.GetGroupInvitation)
		groupInvitations.POST("/:invitationId/accept", routes.AcceptGroupInvitation)
		groupInvitations.POST("/:invitationId/decline", routes.DeclineGroupInvitation)
		groupInvitations.POST("", routes.AddGroupInvitation)
	}

	groupMemberships := router.Group("/memberships").Use(auth.TokenAuth(false))
	{
		groupMemberships.GET("/:username", routes.GetGroupMemberships)
		groupMemberships.POST("", routes.AddGroupMembership)
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
		users.GET("/:username", auth.TokenAuth(true), routes.GetUser)
		users.GET("/:username/invitations", auth.TokenAuth(false), routes.GetGroupInvitationsForUser)
		users.GET("/:username/results", auth.TokenAuth(false), routes.GetResultsForUser)

		users.POST("", routes.RegisterUser)
	}

	secured := router.Group("/secured").Use(auth.TokenAuth(false))
	{
		secured.GET("/ping", routes.Ping)
	}

	router.Run(":8000")
}
