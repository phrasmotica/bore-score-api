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
	GetAllGames(ctx context.Context) (bool, []models.Game)
	GetGame(ctx context.Context, name string) (bool, *models.Game)
	GameExists(ctx context.Context, name string) bool
	AddGame(ctx context.Context, newGame *models.Game) bool
	DeleteGame(ctx context.Context, name string) bool

	GetAllGroups(ctx context.Context) (bool, []models.Group)
	GetGroups(ctx context.Context) (bool, []models.Group)
	GetGroup(ctx context.Context, name string) (RetrieveGroupResult, *models.Group)
	GroupExists(ctx context.Context, name string) bool
	AddGroup(ctx context.Context, newGroup *models.Group) bool
	DeleteGroup(ctx context.Context, name string) bool

	GetAllLinkTypes(ctx context.Context) (bool, []models.LinkType)

	GetAllPlayers(ctx context.Context) (bool, []models.Player)
	GetPlayer(ctx context.Context, username string) (bool, *models.Player)
	PlayerExists(ctx context.Context, username string) bool
	AddPlayer(ctx context.Context, newPlayer *models.Player) bool
	DeletePlayer(ctx context.Context, username string) bool

	GetAllResults(ctx context.Context) (bool, []models.Result)
	AddResult(ctx context.Context, newResult *models.Result) bool
	DeleteResultsWithGame(ctx context.Context, gameName string) (bool, int64)
	ScrubResultsWithPlayer(ctx context.Context, username string) (bool, int64)

	GetAllWinMethods(ctx context.Context) (bool, []models.WinMethod)

	GetSummary(ctx context.Context) (bool, *Summary)
}