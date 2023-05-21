package data

import (
	"context"
	"phrasmotica/bore-score-api/models"
)

type Summary struct {
	GameCount   int64 `json:"gameCount"`
	GroupCount  int64 `json:"groupCount"`
	PlayerCount int64 `json:"playerCount"`
	ResultCount int64 `json:"resultCount"`
}

type IDatabase interface {
	AddApproval(ctx context.Context, newApproval *models.Approval) bool
	GetApprovals(ctx context.Context, resultId string) (bool, []models.Approval)

	GetAllGames(ctx context.Context) (bool, []models.Game)
	GetGame(ctx context.Context, name string) (bool, *models.Game)
	GameExists(ctx context.Context, name string) bool
	AddGame(ctx context.Context, newGame *models.Game) bool
	DeleteGame(ctx context.Context, name string) bool

	GetAllGroups(ctx context.Context) (bool, []models.Group)
	GetGroups(ctx context.Context) (bool, []models.Group)
	GetGroup(ctx context.Context, id string) (bool, *models.Group)
	GetGroupByName(ctx context.Context, name string) (bool, *models.Group)
	GroupExists(ctx context.Context, name string) bool
	AddGroup(ctx context.Context, newGroup *models.Group) bool
	DeleteGroup(ctx context.Context, id string) bool

	GetGroupInvitation(ctx context.Context, invitationId string) (bool, *models.GroupInvitation)
	GetGroupInvitations(ctx context.Context, username string) (bool, []models.GroupInvitation)
	GetGroupInvitationsForGroup(ctx context.Context, groupId string) (bool, []models.GroupInvitation)
	AddGroupInvitation(ctx context.Context, newGroupInvitation *models.GroupInvitation) bool
	UpdateGroupInvitation(ctx context.Context, newGroupInvitation *models.GroupInvitation) bool

	GetGroupMemberships(ctx context.Context, username string) (bool, []models.GroupMembership)
	GetGroupMembershipsForGroup(ctx context.Context, groupId string) (bool, []models.GroupMembership)
	IsInGroup(ctx context.Context, groupId string, username string) bool
	AddGroupMembership(ctx context.Context, newGroupMembership *models.GroupMembership) bool

	GetAllLinkTypes(ctx context.Context) (bool, []models.LinkType)

	GetAllPlayers(ctx context.Context) (bool, []models.Player)
	GetPlayersInGroup(ctx context.Context, groupId string) (bool, []models.Player)
	GetPlayer(ctx context.Context, username string) (bool, *models.Player)
	PlayerExists(ctx context.Context, username string) bool
	AddPlayer(ctx context.Context, newPlayer *models.Player) bool
	DeletePlayer(ctx context.Context, username string) bool

	GetAllResults(ctx context.Context) (bool, []models.Result)
	GetResultsWithPlayer(ctx context.Context, username string) (bool, []models.Result)
	GetResultsForGroup(ctx context.Context, groupId string) (bool, []models.Result)
	GetResult(ctx context.Context, resultId string) (bool, *models.Result)
	ResultExists(ctx context.Context, resultId string) bool
	AddResult(ctx context.Context, newResult *models.Result) bool
	DeleteResultsWithGame(ctx context.Context, gameName string) (bool, int64)
	ScrubResultsWithPlayer(ctx context.Context, username string) (bool, int64)

	GetUser(ctx context.Context, username string) (bool, *models.User)
	GetUserByEmail(ctx context.Context, email string) (bool, *models.User)
	AddUser(ctx context.Context, newUser *models.User) bool
	UserExists(ctx context.Context, username string) bool
	UserExistsByEmail(ctx context.Context, email string) bool

	GetAllWinMethods(ctx context.Context) (bool, []models.WinMethod)

	GetSummary(ctx context.Context) (bool, *Summary)
}
