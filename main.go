package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var winMethods = []WinMethod{
	IndividualScore,
	IndividualWinner,
}

var games = []game{
	{
		ID:       1,
		Name:     "Village Green",
		Synopsis: "A game of pretty gardens and petty grudges.",

		Description: `It's the first day of spring, and there's only one thing on everyone's mind — the Village Green of the Year competition! In just a few months, the judges of this prestigious contest will be visiting, and the village council have finally put you in charge of the preparations. With your newfound authority, you can show those snobs from Lower Aynesmore just what a properly orchestrated floral arrangement looks like!

		In Village Green, you are rival gardeners, tasked by your respective communities with arranging flowers, planting trees, commissioning statues, and building ponds. You must place each element carefully as time is tight and the stakes couldn't be higher! Split your days between acquiring and installing new features for your green and nominating it for one of the competition's many awards. Will your village green become the local laughing stock, or make the neighboring villages green with envy?`,

		MinPlayers: 1,
		MaxPlayers: 5,
		WinMethod:  IndividualScore,
	},
	{
		ID:       2,
		Name:     "Modern Art: The Card Game",
		Synopsis: "Assemble the most valuable art collection.",

		Description: "In Modern Art: The Card Game, the players are art critics, collectors and gallery owners. As it is in art galleries the world over, tastes and opinions change constantly in the world of Modern Art. Today’s treasure is tomorrow’s trash, and no one has more influence on the artists’ values than the players in this game. Which players will exert the most influence on the art market? Who will be the best at anticipating the quickly-changing tastes and opinions of buyers, and thus assemble the highest-valued collection of these new masters? Only the most influential collector will come out on top in Modern Art: The Card Game! Same as Master’s Gallery Bookshelf Game.",

		MinPlayers: 2,
		MaxPlayers: 5,
		WinMethod:  IndividualScore,
	},
	{
		ID:       3,
		Name:     "Love Letter",
		Synopsis: "Can you get a letter to the princess or remove all your rivals? You win either way!",

		Description: `Will your love letter woo the Princess and win her heart? Utilize the characters in the castle to secretly carry your message to the Princess, earning her affection.

		Love Letter is a game of risk, deduction, and luck. Designed by Seiji Kanai, the game features simple rules that create dynamic and fun player interactions. Players attempt to deliver their love letter into the Princess’s hands while keeping other players’ letters away. Powerful cards lead to early gains, but make you a target. Rely on weaker cards for too long and your letter may be tossed in the fire!`,

		MinPlayers: 2,
		MaxPlayers: 4,
		WinMethod:  IndividualScore,
	},
	{
		ID:       4,
		Name:     "Coup",
		Synopsis: "Bluff (and call bluffs!) to victory in this card game with no third chances.",

		Description: `You are head of a family in an Italian city-state, a city run by a weak and corrupt court. You need to manipulate, bluff and bribe your way to power. Your object is to destroy the influence of all the other families, forcing them into exile. Only one family will survive…

		In Coup, you want to be the last player with influence in the game, with influence being represented by face-down character cards in your playing area.`,

		MinPlayers: 2,
		MaxPlayers: 6,
		WinMethod:  IndividualWinner,
	},
}

var players = []player{
	{
		ID:          1,
		Username:    "johannam",
		DisplayName: "Johanna",
	},
	{
		ID:          2,
		Username:    "julianl",
		DisplayName: "Julian",
	},
	{
		ID:          3,
		Username:    "efrimm",
		DisplayName: "Efrim",
	},
	{
		ID:          4,
		Username:    "billyj",
		DisplayName: "Billy",
	},
}

var results = []result{
	{
		ID:        1,
		GameID:    1,
		Timestamp: time.Date(2022, time.January, 22, 10, 34, 0, 0, time.UTC).Unix(),
		Scores: []playerScore{
			{
				PlayerID: 1,
				Score:    25,
			},
			{
				PlayerID: 2,
				Score:    23,
			},
		},
	},
	{
		ID:        2,
		GameID:    1,
		Timestamp: time.Date(2022, time.January, 23, 17, 12, 0, 0, time.UTC).Unix(),
		Scores: []playerScore{
			{
				PlayerID: 1,
				Score:    32,
			},
			{
				PlayerID: 2,
				Score:    34,
			},
		},
	},
	{
		ID:        3,
		GameID:    2,
		Timestamp: time.Date(2022, time.February, 13, 14, 56, 0, 0, time.UTC).Unix(),
		Scores: []playerScore{
			{
				PlayerID: 1,
				Score:    116,
			},
			{
				PlayerID: 2,
				Score:    140,
			},
		},
	},
}

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/games", getGames)
	router.DELETE("/games/:id", deleteGame)

	router.GET("/winMethods", getWinMethods)

	router.GET("/players", getPlayers)
	router.POST("/players", postPlayer)
	router.DELETE("/players/:username", deletePlayer)

	router.GET("/results", getResults)
	router.POST("/results", postResult)

	router.Run("localhost:8000")
}

func getGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, games)
}

func deleteGame(c *gin.Context) {
	gameId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid game id"})
		return
	}

	if !gameExists(games, gameId) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %d does not exist", gameId)})
		return
	}

	games = removeGame(games, gameId)
	c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func getWinMethods(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, winMethods)
}

func getPlayers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, players)
}

func postPlayer(c *gin.Context) {
	var newPlayer player

	if err := c.BindJSON(&newPlayer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if playerExistsByUsername(players, newPlayer.Username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s already exists", newPlayer.Username)})
		return
	}

	newPlayer.ID = getMaxPlayerId(players) + 1

	players = append(players, newPlayer)
	c.IndentedJSON(http.StatusCreated, newPlayer)
}

func deletePlayer(c *gin.Context) {
	username := c.Param("username")

	if !playerExistsByUsername(players, username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %s does not exist", username)})
		return
	}

	players = removePlayer(players, username)
	c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func getResults(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, results)
}

func postResult(c *gin.Context) {
	var newResult result

	if err := c.BindJSON(&newResult); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body format"})
		return
	}

	if !gameExists(games, newResult.GameID) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("game %d does not exist", newResult.GameID)})
		return
	}

	for _, score := range newResult.Scores {
		if !playerExists(players, score.PlayerID) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("player %d does not exist", score.PlayerID)})
			return
		}
	}

	newResult.ID = getMaxResultId(results) + 1

	results = append(results, newResult)
	c.IndentedJSON(http.StatusCreated, newResult)
}
