package main

type WinMethod string

const (
	IndividualScore  WinMethod = "Individual Score"
	IndividualWinner WinMethod = "Individual Winner"
)

type LinkType string

const (
	OfficialWebsite LinkType = "Official Website"
	BoardGameGeek   LinkType = "BoardGameGeek"
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
	Links       []Link    `json:"links"`
}

type Link struct {
	Type LinkType `json:"type"`
	Link string   `json:"link"`
}

type player struct {
	ID          int    `json:"id"`
	TimeCreated int64  `json:"timeCreated"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type playerScore struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
	IsWinner bool   `json:"isWinner"`
}

type result struct {
	ID        int           `json:"id"`
	GameID    int           `json:"gameId"`
	Timestamp int64         `json:"timestamp"`
	Scores    []playerScore `json:"scores"`
}
