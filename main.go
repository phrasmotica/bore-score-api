package main

import (
	"phrasmotica/bore-score-api/auth"
	"phrasmotica/bore-score-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(auth.CORSMiddleware())

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
		groups.DELETE("/:name", routes.DeleteGroup) // TODO: select by group ID instead of name
	}

	groupInvitations := router.Group("/invitations").Use(auth.TokenAuth(false))
	{
		groupInvitations.GET("/:invitationId", routes.GetGroupInvitation)
		groupInvitations.POST("/:invitationId/accept", routes.AcceptGroupInvitation)
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
		users.GET("/:username/invitations", auth.TokenAuth(false), routes.GetGroupInvitations)
		users.POST("", routes.RegisterUser)
	}

	secured := router.Group("/secured").Use(auth.TokenAuth(false))
	{
		secured.GET("/ping", routes.Ping)
	}

	router.Run(":8000")
}
