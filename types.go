package main

type game struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
