package models

type WinMethod string

const (
	IndividualScore  WinMethod = "Individual Score"
	IndividualWinner WinMethod = "Individual Winner"
)

var WinMethods = []WinMethod{
	IndividualScore,
	IndividualWinner,
}

type LinkType string

const (
	OfficialWebsite LinkType = "Official Website"
	BoardGameGeek   LinkType = "BoardGameGeek"
)

var LinkTypes = []LinkType{
	OfficialWebsite,
	BoardGameGeek,
}

type Game struct {
	ID          string    `json:"id" bson:"id"`
	Name        string    `json:"name" bson:"name"`
	TimeCreated int64     `json:"timeCreated" bson:"timeCreated"`
	DisplayName string    `json:"displayName" bson:"displayName"`
	Synopsis    string    `json:"synopsis" bson:"synopsis"`
	Description string    `json:"description" bson:"description"`
	MinPlayers  int       `json:"minPlayers" bson:"minPlayers"`
	MaxPlayers  int       `json:"maxPlayers" bson:"maxPlayers"`
	WinMethod   WinMethod `json:"winMethod" bson:"winMethod"`
	Links       []Link    `json:"links" bson:"links"`
}

type Link struct {
	Type LinkType `json:"type" bson:"type"`
	Link string   `json:"link" bson:"link"`
}

type Player struct {
	ID          string `json:"id" bson:"id"`
	Username    string `json:"username" bson:"username"`
	TimeCreated int64  `json:"timeCreated" bson:"timeCreated"`
	DisplayName string `json:"displayName" bson:"displayName"`
}

type PlayerScore struct {
	Username string `json:"username" bson:"username"`
	Score    int    `json:"score" bson:"score"`
	IsWinner bool   `json:"isWinner" bson:"isWinner"`
}

type Result struct {
	ID        string        `json:"id" bson:"id"`
	GameName  string        `json:"gameName" bson:"gameName"`
	Timestamp int64         `json:"timestamp" bson:"timestamp"`
	Scores    []PlayerScore `json:"scores" bson:"scores"`
}
