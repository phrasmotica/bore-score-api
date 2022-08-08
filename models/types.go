package models

type Game struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	TimeCreated int64  `json:"timeCreated" bson:"timeCreated"`
	DisplayName string `json:"displayName" bson:"displayName"`
	Synopsis    string `json:"synopsis" bson:"synopsis"`
	Description string `json:"description" bson:"description"`
	MinPlayers  int    `json:"minPlayers" bson:"minPlayers"`
	MaxPlayers  int    `json:"maxPlayers" bson:"maxPlayers"`
	WinMethod   string `json:"winMethod" bson:"winMethod"`
	ImageLink   string `json:"imageLink" bson:"imageLink"`
	Links       []Link `json:"links" bson:"links"`
}

type Link struct {
	Type LinkTypeName `json:"type" bson:"type"`
	Link string       `json:"link" bson:"link"`
}

type LinkType struct {
	ID          string       `json:"id" bson:"id"`
	Name        LinkTypeName `json:"name" bson:"name"`
	TimeCreated int64        `json:"timeCreated" bson:"timeCreated"`
	DisplayName string       `json:"displayName" bson:"displayName"`
}

type LinkTypeName string

const (
	OfficialWebsite LinkTypeName = "official-website"
	BoardGameGeek   LinkTypeName = "board-game-geek"
)

type Player struct {
	ID             string `json:"id" bson:"id"`
	Username       string `json:"username" bson:"username"`
	TimeCreated    int64  `json:"timeCreated" bson:"timeCreated"`
	DisplayName    string `json:"displayName" bson:"displayName"`
	ProfilePicture string `json:"profilePicture" bson:"profilePicture"`
}

type PlayerScore struct {
	Username string `json:"username" bson:"username"`
	Score    int    `json:"score" bson:"score"`
	IsWinner bool   `json:"isWinner" bson:"isWinner"`
}

type Result struct {
	ID               string        `json:"id" bson:"id"`
	GameName         string        `json:"gameName" bson:"gameName"`
	GroupName        string        `json:"groupName" bson:"groupName"`
	TimeCreated      int64         `json:"timeCreated" bson:"timeCreated"`
	TimePlayed       int64         `json:"timePlayed" bson:"timePlayed"`
	Notes            string        `json:"notes" bson:"notes"`
	CooperativeScore int           `json:"cooperativeScore" bson:"cooperativeScore"`
	CooperativeWin   bool          `json:"cooperativeWin" bson:"cooperativeWin"`
	Scores           []PlayerScore `json:"scores" bson:"scores"`
}

type WinMethod struct {
	ID          string        `json:"id" bson:"id"`
	Name        WinMethodName `json:"name" bson:"name"`
	TimeCreated int64         `json:"timeCreated" bson:"timeCreated"`
	DisplayName string        `json:"displayName" bson:"displayName"`
}

type WinMethodName string

const (
	IndividualScore  WinMethodName = "individual-score"
	IndividualWin    WinMethodName = "individual-win"
	CooperativeScore WinMethodName = "cooperative-score"
	CooperativeWin   WinMethodName = "cooperative-win"
)
