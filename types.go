package main

type WinMethod string

const (
	IndividualScore  WinMethod = "Individual Score"
	IndividualWinner WinMethod = "Individual Winner"
)

type game struct {
	ID          int       `json:"id"`
	TimeCreated int64     `json:"timeCreated"`
	Name        string    `json:"name"`
	Synopsis    string    `json:"synopsis"`
	Description string    `json:"description"`
	MinPlayers  int       `json:"minPlayers"`
	MaxPlayers  int       `json:"maxPlayers"`
	WinMethod   WinMethod `json:"winMethod"`
}

type player struct {
	ID          int    `json:"id"`
	TimeCreated int64  `json:"timeCreated"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type playerScore struct {
	PlayerID int  `json:"playerId"`
	Score    int  `json:"score"`
	IsWinner bool `json:"isWinner"`
}

type result struct {
	ID        int           `json:"id"`
	GameID    int           `json:"gameId"`
	Timestamp int64         `json:"timestamp"`
	Scores    []playerScore `json:"scores"`
}
