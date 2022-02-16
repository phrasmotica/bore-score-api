package main

type GameType string

const (
	Score GameType = "Score"
)

type game struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	GameType   GameType `json:"gameType"`
	MinPlayers int      `json:"minPlayers"`
	MaxPlayers int      `json:"maxPlayers"`
}

type player struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type playerScore struct {
	PlayerID int `json:"playerId"`
	Score    int `json:"score"`
}

type result struct {
	ID     int           `json:"id"`
	GameID int           `json:"gameId"`
	Scores []playerScore `json:"scores"`
}
