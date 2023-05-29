package routes

import (
	"context"
	"net/http"
	"phrasmotica/bore-score-api/models"

	"github.com/gin-gonic/gin"
)

type LeaderboardResponse struct {
	GroupID     string `json:"groupId" bson:"groupId"`
	GameID      string `json:"gameId" bson:"gameId"`
	Leaderboard []Rank `json:"leaderboard" bson:"leaderboard"`
}

type Rank struct {
	// TODO: add number of games played, wins/draws/losses/etc
	Username     string `json:"username" bson:"username"`
	PointsScored int    `json:"pointsScored" bson:"pointsScored"`
}

func GetLeaderboard(c *gin.Context) {
	groupId := c.Param("groupId")
	gameId := c.Param("gameId")

	callingUsername := c.GetString("username")

	ctx := context.TODO()

	success, group := db.GetGroup(ctx, groupId)
	if !success {
		Error.Printf("Group %s not found\n", groupId)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	success, game := db.GetGame(ctx, gameId)
	if !success {
		Error.Printf("Game %s not found\n", gameId)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if !db.IsInGroup(ctx, group.ID, callingUsername) {
		Error.Printf("User %s is not in group %s\n", callingUsername, group.ID)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	success, response := computeLeaderboard(ctx, group, game)
	if !success {
		Error.Printf("Failed to compute leaderboard for game %s in group %s\n", game.ID, group.ID)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	Info.Printf("Computed leaderboard for game %s in group %s\n", game.ID, group.ID)

	c.IndentedJSON(http.StatusOK, response)
}

func computeLeaderboard(ctx context.Context, group *models.Group, game *models.Game) (bool, *LeaderboardResponse) {
	success, results := db.GetResultsForGroupAndGame(ctx, group.ID, game.Name)
	if !success {
		return false, nil
	}

	pointsMap := map[string]int{}

	for _, r := range results {
		for _, s := range r.Scores {
			_, exists := pointsMap[s.Username]
			if exists {
				pointsMap[s.Username] += s.Score
			} else {
				pointsMap[s.Username] = s.Score
			}
		}
	}

	leaderboard := []Rank{}

	for player, score := range pointsMap {
		leaderboard = append(leaderboard, Rank{
			Username:     player,
			PointsScored: score,
		})
	}

	return true, &LeaderboardResponse{
		GroupID:     group.ID,
		GameID:      game.ID,
		Leaderboard: leaderboard,
	}
}
